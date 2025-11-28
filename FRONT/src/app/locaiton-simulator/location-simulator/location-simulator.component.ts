import { Component, AfterViewInit, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as L from 'leaflet';
import { SimulatorService } from '../../shared/services/simulator-service.service';
import { Location } from '../../models/location.model';

@Component({
  selector: 'app-location-simulator',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './location-simulator.component.html',
  styleUrls: ['./location-simulator.component.css']
})
export class LocationSimulatorComponent implements AfterViewInit {

  private map!: L.Map;
  private marker?: L.Marker;

  // Poruka za korisnika
  public message: string = '';
  public messageType: 'success' | 'info' = 'info';

  // --- DEFINICIJA IKONICE ---
  // Definisemo kako nas "cikica" izgleda
  private personIcon = L.icon({
    iconUrl: 'assets/gps-marker.png', // <--- OVO MORA BITI PUTANJA DO TVOJE SLIKE U ASSETS
    // Ako nemas sliku jos uvek, za test stavi ovaj URL:
    // iconUrl: 'https://cdn-icons-png.flaticon.com/512/9131/9131546.png',

    iconSize: [40, 40],   // Velicina ikonice u pikselima [sirina, visina]
    iconAnchor: [20, 40], // Tacka ikonice koja "bode" lokaciju (polovina sirine, puna visina)
    popupAnchor: [0, -40] // Gde se otvara popup u odnosu na ikonicu
  });

  constructor(private simulatorService: SimulatorService) { }

  ngAfterViewInit(): void {
    this.initMap();
  }

  private initMap(): void {
    // Default centar (npr. Beograd) ako nema sacuvane lokacije
    let initialLat = 44.8125;
    let initialLng = 20.4612;
    let zoomLevel = 13;

    // Proveri da li imamo sacuvanu lokaciju
    const savedLocation = this.simulatorService.getLocation();
    if (savedLocation) {
      initialLat = savedLocation.latitude;
      initialLng = savedLocation.longitude;
      zoomLevel = 15; // Malo blizi zoom ako imamo tacnu lokaciju
    }

    this.map = L.map('mapSimulator').setView([initialLat, initialLng], zoomLevel);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // Ako imamo sacuvanu lokaciju, odmah dodaj marker
    if (savedLocation) {
      this.setMarker(savedLocation.latitude, savedLocation.longitude);
      this.message = "Loaded saved location.";
    } else {
      this.message = "Click on the map to set your location.";
    }

    // Slusaj klik na mapu
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      this.setMarker(e.latlng.lat, e.latlng.lng);
      this.message = "Location selected. Click 'Confirm' to save.";
      this.messageType = 'info';
    });
  }

  private setMarker(lat: number, lng: number): void {
    // Ako marker vec postoji, ukloni ga (zelimo samo jedan)
    if (this.marker) {
      this.map.removeLayer(this.marker);
    }

    this.marker = L.marker([lat, lng], { icon: this.personIcon }).addTo(this.map);
  }

  public confirmLocation(): void {
    if (!this.marker) {
      alert("Please select a location on the map first.");
      return;
    }

    const latLng = this.marker.getLatLng();
    const location: Location = {
      latitude: latLng.lat,
      longitude: latLng.lng
    };

    this.simulatorService.setLocation(location);

    this.message = "Location saved successfully! Tours will now use this position.";
    this.messageType = 'success';
  }
}