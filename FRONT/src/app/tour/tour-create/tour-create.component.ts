import { Component, AfterViewInit, ViewChild, ElementRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from "@angular/common";
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Router } from '@angular/router';
import * as L from 'leaflet';

// --- MODELS ---

export interface Keypoint {
  id?: string;
  tourId: string;
  title: string;
  latitude: number;
  longitude: number;
  description: string;
  imageUrl: string;
}

export enum RouteMode {
  Walking = 'WALKING',
  Cycling = 'CYCLING',
  Driving = 'DRIVING'
}

export interface RouteOption {
  mode: RouteMode;
  duration: number; // Nanosekunde za backend
}

export interface RouteOptionInput {
  mode: RouteMode;
  hours: number;
  minutes: number;
}

export interface Tour {
  title: string;
  description: string;
  difficulty: string;
  tags: string; // Ovo ostaje string za backend
  status: string;
  price: number;
  distance: number; 
  keypoints: Keypoint[];
  route_options: RouteOption[]; 
}

@Component({
  selector: 'app-tour-create',
  standalone: true,
  imports: [FormsModule, CommonModule, HttpClientModule],
  templateUrl: './tour-create.component.html',
  styleUrls: ['./tour-create.component.css']
})
export class TourCreateComponent implements AfterViewInit {
  
  // Map logic
  public map!: L.Map;
  public pathLine?: L.Polyline;
  public markers: L.Marker[] = [];
  public keypoints: Keypoint[] = [];
  public currentMarker?: L.Marker;
  public files: File[] = [];

  // State logic
  public editing: boolean = false;
  public currentStep: number = 1; 
  
  @ViewChild('fileInput') fileInput!: ElementRef<HTMLInputElement>;

  // Data holders
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
    distance: 0,
    keypoints: [],
    route_options: []
  };

  // --- TAGS LOGIC (NOVO) ---
  public availableTags: string[] = [
    'Nature', 
    'History', 
    'Culture', 
    'Food & Drink', 
    'Adventure', 
    'Relaxation', 
    'Urban', 
    'Nightlife',
    'Museums',
    'Riverside'
  ];
  
  // Ovde cuvamo sta je korisnik kliknuo pre nego sto pretvorimo u string
  public selectedTags: string[] = [];

  // Route Option logic
  public availableModes = [RouteMode.Walking, RouteMode.Cycling, RouteMode.Driving];
  public currentRouteOptionInput: RouteOptionInput = {
    mode: RouteMode.Walking,
    hours: 0,
    minutes: 0
  };
  public addedRouteOptions: RouteOptionInput[] = [];

  // Modals
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
      if (!this.editing && this.currentStep === 1) {
        this.addMarker(e.latlng.lat, e.latlng.lng);
      }
    });
  }

  // --- STEP 1 LOGIC (Keypoints) ---

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
      if (this.currentStep === 1) {
        this.selectMarker(marker, keypoint);
      }
    });

    this.selectMarker(marker, keypoint);
  }

  private selectMarker(marker: L.Marker, keypoint: Keypoint): void {
    this.currentMarker = marker;
    this.currentKeypoint = keypoint; 
    this.editing = true; 
  }

  private resetForm(): void {
    this.currentMarker = undefined;
    this.currentKeypoint = {
      tourId: '',
      title: '',
      description: '',
      latitude: 0,
      longitude: 0,
      imageUrl: ''
    };
    this.editing = false; 
    if(this.fileInput) this.fileInput.nativeElement.value = '';
  }

  public savePoint(): void {
    if(!this.currentKeypoint.title) {
        alert("Please enter a title for this keypoint.");
        return;
    }
    this.resetForm();
  }

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
    this.resetForm(); 
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

  // --- NAVIGATION LOGIC ---

  public goToStep2(): void {
    if(this.keypoints.length < 2) {
        alert("You need at least 2 keypoints to create a tour.");
        return;
    }
    if (this.editing) {
       alert("Please save or discard the current point before finishing.");
       return;
    }
    this.calculateTourDistance();
    this.currentStep = 2;
    this.editing = false;
  }

  public backToStep1(): void {
    this.currentStep = 1;
  }

  public goToStep3(): void {
    if (this.addedRouteOptions.length === 0) {
        alert("Please add at least one route option (e.g., Walking).");
        return;
    }
    this.currentStep = 3;
  }

  public backToStep2(): void {
    this.currentStep = 2;
  }

  // --- STEP 2 LOGIC (Distance & RouteOptions) ---

  private calculateTourDistance(): void {
    let totalDistanceMeters = 0;
    for (let i = 0; i < this.keypoints.length - 1; i++) {
      const point1 = L.latLng(this.keypoints[i].latitude, this.keypoints[i].longitude);
      const point2 = L.latLng(this.keypoints[i+1].latitude, this.keypoints[i+1].longitude);
      totalDistanceMeters += point1.distanceTo(point2);
    }
    this.currentTour.distance = parseFloat((totalDistanceMeters / 1000).toFixed(2));
  }

  public addRouteOption(): void {
    const input = this.currentRouteOptionInput;
    if (input.hours === 0 && input.minutes === 0) {
        alert("Duration cannot be zero.");
        return;
    }
    const exists = this.addedRouteOptions.find(o => o.mode === input.mode);
    if (exists) {
        alert(`Option for ${input.mode} already exists. Please remove it first if you want to change it.`);
        return;
    }
    this.addedRouteOptions.push({ ...input });
    this.currentRouteOptionInput = { mode: RouteMode.Walking, hours: 0, minutes: 0 };
  }

  public removeRouteOption(index: number): void {
    this.addedRouteOptions.splice(index, 1);
  }

  public formatDuration(opt: RouteOptionInput): string {
    const h = opt.hours > 0 ? `${opt.hours}h ` : '';
    const m = opt.minutes > 0 ? `${opt.minutes}min` : '';
    return (h + m).trim();
  }

  // --- TAGS LOGIC (METODE) ---
  
  public toggleTag(tag: string): void {
    if (this.selectedTags.includes(tag)) {
      // Ako je vec selektovan, izbaci ga
      this.selectedTags = this.selectedTags.filter(t => t !== tag);
    } else {
      // Ako nije, dodaj ga
      this.selectedTags.push(tag);
    }
  }

  // --- STEP 3 LOGIC (Final & Save) ---

  public saveTour(): void {
    this.currentTour.keypoints = this.keypoints;
    this.currentTour.status = "Published"; 

    // Konverzija Tagova: Niz stringova -> Jedan string odvojen zarezima
    this.currentTour.tags = this.selectedTags.join(', ');

    // Konverzija Vremena
    this.currentTour.route_options = this.addedRouteOptions.map(opt => {
        const totalMinutes = (opt.hours * 60) + opt.minutes;
        const durationNanoseconds = totalMinutes * 60 * 1_000_000_000; 
        return {
            mode: opt.mode,
            duration: durationNanoseconds
        };
    });

    if(!this.currentTour.title || !this.currentTour.description || this.currentTour.price <= 0) {
        alert("Please fill in all Tour Details fields correctly.");
        return;
    }
    if(this.selectedTags.length === 0) {
        alert("Please select at least one tag.");
        return;
    }

    this.http.post('http://localhost:8080/tour', this.currentTour).subscribe({
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
    this.router.navigate(['/tour-list']); 
  }

  public closeErrorModal() {
    this.showErrorModal = false;
  }
}