import { Component, AfterViewInit, ViewChild, ElementRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from "@angular/common";
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Router } from '@angular/router';
import * as L from 'leaflet';
import { Tour } from '../model/tour';
import { Keypoint } from '../model/keypoint';

@Component({
  selector: 'app-tour-create',
  standalone: true,
  imports: [FormsModule, CommonModule, HttpClientModule],
  templateUrl: './tour-create.component.html',
  styleUrls: ['./tour-create.component.css']
})
export class TourCreateComponent implements AfterViewInit {
  
  public map!: L.Map;
  public pathLine?: L.Polyline;
  public markers: L.Marker[] = [];
  public keypoints: Keypoint[] = [];
  public currentMarker?: L.Marker;
  public files: File[] = [];

  public editing: boolean = false;
  public keyPointMode: boolean = true;
  
  @ViewChild('fileInput') fileInput!: ElementRef<HTMLInputElement>;

  public currentKeypoint: Keypoint = {
    tourId: '',
    title: '',
    description: '',
    latitude: 0,
    longitude: 0,
    imageUrl: ''
  };

  public currentTour: Tour = {
    title: '',
    description: '',
    difficulty: 'Medium',
    status: 'Draft',
    tags: '',
    price: 0,
    keypoints: []
  };

  public showSuccessModal = false;
  public showErrorModal = false;

  constructor(private http: HttpClient, private router: Router) {}

  ngAfterViewInit(): void {
    this.initMap();
  }

  private initMap(): void {
    const lat = 44.8125; 
    const lng = 20.4612;

    this.map = L.map('map').setView([lat, lng], 13);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    this.map.on('click', (e: L.LeafletMouseEvent) => {
      // Dozvoli dodavanje samo ako NE editujemo trenutno tacku
      if (!this.editing && this.keyPointMode) {
        this.addMarker(e.latlng.lat, e.latlng.lng);
      }
    });
  }

  private addMarker(lat: number, lng: number): void {
    const marker = L.marker([lat, lng]).addTo(this.map);

    const keypoint: Keypoint = { 
      tourId: '',
      latitude: lat, 
      longitude: lng, 
      title: '', 
      description: '', 
      imageUrl: '' 
    };

    this.markers.push(marker);
    this.keypoints.push(keypoint);
    this.files.push(null as any); 

    this.redrawPath();

    marker.on('click', () => {
      if (this.keyPointMode) {
        this.selectMarker(marker, keypoint);
      }
    });

    // Odmah ulazimo u edit mod za novu tacku
    this.selectMarker(marker, keypoint);
  }

  private selectMarker(marker: L.Marker, keypoint: Keypoint): void {
    this.currentMarker = marker;
    this.currentKeypoint = keypoint; // Referenca na objekat u nizu
    this.editing = true; // Blokiramo mapu dok se ovo ne zavrsi
  }

  // --- NOVA FUNKCIJA: Samo resetuje formu, bez validacije ---
  private resetForm(): void {
    this.currentMarker = undefined;
    // Resetujemo currentKeypoint na novi prazan objekat da ne bi ostala referenca
    this.currentKeypoint = {
      tourId: '',
      title: '',
      description: '',
      latitude: 0,
      longitude: 0,
      imageUrl: ''
    };
    this.editing = false; // Odblokiramo mapu
    if(this.fileInput) this.fileInput.nativeElement.value = '';
  }

  // Klik na "Save Point"
  public savePoint(): void {
    if(!this.currentKeypoint.title) {
        alert("Please enter a title for this keypoint.");
        return;
    }
    // Ako je validno, samo resetujemo formu (podaci su vec u nizu keypoints jer koristimo referencu)
    this.resetForm();
  }

  // Klik na "Discard Point"
  public discardMarker(): void {
    if (!this.currentMarker) return;

    const marker = this.currentMarker;
    const index = this.markers.indexOf(marker);

    if (index !== -1) {
      this.map.removeLayer(marker);
      this.markers.splice(index, 1);
      this.keypoints.splice(index, 1);
      this.files.splice(index, 1);
    }

    this.redrawPath();
    this.resetForm(); // Resetujemo bez pitanja za Title
  }

  public finishKeypoints(): void {
    if(this.keypoints.length < 2) {
        alert("You need at least 2 keypoints to create a tour.");
        return;
    }
    // Ako je ostala neka tacka nezavrsena (edit mode), sacuvaj je ili odbaci
    if (this.editing) {
       alert("Please save or discard the current point before finishing.");
       return;
    }

    this.keyPointMode = false;
    this.editing = false;
  }

  public backToKeypoints(): void {
    this.keyPointMode = true;
    this.editing = false;
  }

  private redrawPath(): void {
    const points = this.keypoints.map(k => [k.latitude, k.longitude]) as [number, number][];

    if (this.pathLine) {
      this.map.removeLayer(this.pathLine);
    }

    if (points.length >= 2) {
      this.pathLine = L.polyline(points, { color: '#0b76ef', weight: 4 }).addTo(this.map);
    }
  }   

  public onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    
    const file = input.files[0];
    const previewUrl = URL.createObjectURL(file);
    
    const popupContent = `<div style="text-align:center;"><img src="${previewUrl}" width="100" style="border-radius:4px;"></div>`;

    if (this.currentMarker) {
      this.currentMarker.bindTooltip(popupContent, {
        direction: 'top',
        opacity: 1,
        className: 'hover-popup'
      }).openTooltip();
    }

    this.currentKeypoint.imageUrl = file.name;
    
    const index = this.keypoints.indexOf(this.currentKeypoint);
    if (index !== -1) {
        this.files[index] = file;
    }
  }

  public saveTour(): void {
    this.currentTour.keypoints = this.keypoints;
    this.currentTour.status = "Published"; 

    if(!this.currentTour.title || !this.currentTour.description || this.currentTour.price <= 0) {
        alert("Please fill in all Tour Details fields correctly.");
        return;
    }

    this.http.post('http://localhost:8083/tour', this.currentTour).subscribe({
      next: (response) => {
        this.showSuccessModal = true;
      },
      error: (error) => {
        console.error("Error creating tour:", error);
        this.showErrorModal = true;
      }
    });
  }

  public closeSuccessModal() {
    this.showSuccessModal = false;
    this.router.navigate(['/tours']); 
  }

  public closeErrorModal() {
    this.showErrorModal = false;
  }
}