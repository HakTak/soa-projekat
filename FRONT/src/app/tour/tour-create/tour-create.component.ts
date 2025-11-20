import { Component, AfterViewInit, ViewChild, ElementRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Keypoint } from '../model/keypoint';
import { Tour } from '../model/tour';
import {CommonModule} from "@angular/common";
import * as L from 'leaflet';

@Component({
  selector: 'app-tour-create',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './tour-create.component.html',
  styleUrl: './tour-create.component.css'
})
export class TourCreateComponent {
  public map!: L.Map;
  public pathLine?: L.Polyline;
  public markers: L.Marker[] = [];       // Leaflet markers
  public keypoints: Keypoint[] = [];      // corresponding data
  public currentMarker?: L.Marker;        // the marker being edited
  public files: File[] = [];
  public currentKeypoint: Keypoint = {
  title: '',
  description: '',
  latitude: 0,
  longitude: 0,
  imageUrl: ''
}
public currentTour: Tour = {
  title: '',
  description: '',
  difficulty: '',
  status: 'Draft',
  tags: '',
  price: 0,
  keypoints: []
}
  public editing: boolean = false;
  public keyPointMode: boolean = true;
  @ViewChild('fileInput') fileInput!: ElementRef<HTMLInputElement>;

  ngAfterViewInit(): void {
    this.initMap();
  }

  private initMap(): void {
    const lat = 44.7866;
    const lng = 20.4489;

    this.map = L.map('map').setView([lat, lng], 13);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);
     this.map.on('click', (e: L.LeafletMouseEvent) => {
        if (!this.editing) {
          this.addMarker(e.latlng.lat, e.latlng.lng);
        }
     });
  }

  private addMarker(lat: number, lng: number): void {
    const marker = L.marker([lat, lng]).addTo(this.map);

    const keypoint: Keypoint = { 
      latitude: lat, 
      longitude: lng, 
      title: '', 
      description: '', 
      imageUrl: ''
    };

    this.markers.push(marker);
    this.keypoints.push(keypoint);

    this.redrawPath();

    marker.on('click', () => {
      this.selectMarker(marker, keypoint);
      this.keyPointMode = true;
    });
    this.keyPointMode = true;
    this.selectMarker(marker, keypoint);
  }

  private selectMarker(marker: L.Marker, keypoint: Keypoint): void {
    this.currentMarker = marker;
    this.currentKeypoint = keypoint;
    this.editing = true;
  }

  public nextMarker(): void {
    this.currentMarker = undefined;
    this.currentKeypoint = {
      title: '',
      description: '',
      latitude: 0,
      longitude: 0,
      imageUrl: ''
    }
    this.editing = false;
    this.fileInput.nativeElement.value = '';
  }

  public discardMarker(): void {
    const marker = this.currentMarker!;
    const keypoint = this.currentKeypoint!;

    const index = this.markers.indexOf(marker);

    if (index !== -1) {
      this.map.removeLayer(marker);
      this.markers.splice(index, 1);
      this.keypoints.splice(index, 1);
      this.files.splice(index, 1);
    }

    this.currentMarker = undefined;
    this.currentKeypoint = {
      title: '',
      description: '',
      latitude: 0,
      longitude: 0,
      imageUrl: ''
    }
    this.editing = false;
    this.fileInput.nativeElement.value = '';
    this.redrawPath();
  }

  public finish(): void {
    this.nextMarker();
    this.keyPointMode = false;
  }

  private redrawPath(): void {
    const points = this.keypoints.map(k => [k.latitude, k.longitude]) as [number, number][];

    // Remove old line
    if (this.pathLine) {
      this.map.removeLayer(this.pathLine);
    }

    if (points.length >= 2) {
      this.pathLine = L.polyline(points, { color: 'blue', weight: 4 }).addTo(this.map);
    }
  }   

  public onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const previewUrl = URL.createObjectURL(input.files[0]);
    const popupContent = `
        <div style="text-align:center;">
          <img src="${previewUrl}" width="120">
        </div>
      `;

      if (this.currentMarker)
        this.currentMarker.bindTooltip(popupContent, {
          direction: 'top',
          opacity: 0.95,
          className: 'hover-popup'
        });
        this.files.push(input.files[0]);
    }

  public saveTour(): void {
    const formData = new FormData();
    this.currentTour.keypoints = this.keypoints;
    formData.append('tour', JSON.stringify(this.currentTour));
    this.keypoints.forEach((_,index) => {
      if (this.files[index]) {
          formData.append('images', this.files[index], `image_${index}.jpg`);
      }
    });
    this.markers.forEach(marker => this.map.removeLayer(marker));
    this.markers = [];
    this.keypoints = [];
    this.files = [];
    this.redrawPath();
    this.keyPointMode = true;
    // Here you would typically send formData to your backend via HTTP
  }
}
