import { Component, AfterViewInit, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as L from 'leaflet';
import { SimulatorService } from '../../shared/services/simulator-service.service';
import { Location } from '../../models/location.model';
import { TourExecutionStatus, CompletedKeypoint, TouristPosition,  TourExecution } from '../model/tour-execution';
import { ActivatedRoute } from '@angular/router';
import { Tour } from '../model/tour'
import { HttpClient} from '@angular/common/http';


@Component({
  selector: 'app-tour-execution',
  standalone: true,
  imports: [],
  templateUrl: './tour-execution.component.html',
  styleUrl: './tour-execution.component.css'
})
export class TourExecutionComponent {
  private map!: L.Map;
  private userMarker?: L.Marker;
  private savedLocation: Location | null = null;
  public tourExecution: TourExecution | null = null;
  public activatedTour: Tour | null = null;
  private intervalId: any;

  private personIcon = L.icon({
    iconUrl: 'assets/gps-marker.png',
    iconSize: [40, 40],
    iconAnchor: [20, 40],
    popupAnchor: [0, -40]
  });

  private keypointIcon = L.icon({
    iconUrl: 'assets/keypoint-marker.png',
    iconSize: [30, 30],
    iconAnchor: [15, 30],
    popupAnchor: [0, -30]
  });

  private keypointMarkers: L.Marker[] = [];
  constructor(
    private simulatorService: SimulatorService,
    private route: ActivatedRoute,
    private http: HttpClient
  ) {
    this.savedLocation = this.simulatorService.getLocation();

    this.route.paramMap.subscribe(params => {
      const tourId = params.get('id');
      if (tourId) {
        this.fetchTour(tourId);
      }
    });
  }

  ngAfterViewInit(): void {
    this.initMap();
  }

  ngOnDestroy(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  private initMap(): void {
    const initialLat = this.savedLocation?.latitude || 44.8125;
    const initialLng = this.savedLocation?.longitude || 20.4612;
    const zoomLevel = this.savedLocation ? 15 : 13;

    this.map = L.map('tourExecutionMap').setView([initialLat, initialLng], zoomLevel);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // Add user marker if saved location exists
    if (this.savedLocation) {
      this.userMarker = L.marker([this.savedLocation.latitude, this.savedLocation.longitude], { icon: this.personIcon })
        .addTo(this.map);
    }

    // Update user marker every 10 seconds
    this.intervalId = setInterval(() => this.updateMarker(), 10000);
  }

  private updateMarker(): void {
  const newLocation = this.simulatorService.getLocation();
  if (!newLocation || !this.tourExecution || !this.activatedTour) return;

  this.savedLocation = newLocation;

  // Move the user marker
  if (this.userMarker) {
    this.userMarker.setLatLng([this.savedLocation!.latitude, this.savedLocation!.longitude]);
  } else {
    this.userMarker = L.marker([this.savedLocation!.latitude, this.savedLocation!.longitude], { icon: this.personIcon })
      .addTo(this.map);
  }

  // Check proximity to next uncompleted keypoint
  const completedIds = new Set(this.tourExecution.completedKeyPoints.map(kp => kp.keyPointId));
  const proximityThresholdMeters = 20; // e.g., consider "close" if within 20 meters

  this.activatedTour.keypoints
    .filter(kp => kp.id !== undefined && !completedIds.has(kp.id))
    .forEach(kp => {
      const distance = this.getDistanceMeters(
        this.savedLocation!.latitude,
        this.savedLocation!.longitude,
        kp.latitude,
        kp.longitude
      );

      if (distance <= proximityThresholdMeters) {
        console.log(`Keypoint ${kp.id} completed!`);
        this.tourExecution!.completedKeyPoints.push({
          keyPointId: kp.id!,
          completedAt: new Date()
        });
      }
    });

  // Redraw keypoints (uncompleted ones only)
  this.drawKeypoints();
  this.tourExecution.lastActivityAt = new Date()
  this.tourExecution.currentPosition = {
            lat: this.savedLocation!.latitude,
            lng: this.savedLocation!.longitude
          }

  // Check if all keypoints are now completed
  const totalKeypoints = this.activatedTour.keypoints.length;
  const completedCount = this.tourExecution.completedKeyPoints.length;

  if (completedCount >= totalKeypoints) {
    console.log('All keypoints completed! Tour finished.');
    this.tourExecution.status = TourExecutionStatus.COMPLETED;    
    this.tourExecution.finishedAt = new Date()  

    //http here
  }  
  // Optionally recenter map on user
  this.map.setView([newLocation.latitude, newLocation.longitude]);

  console.log('Updated user position:', newLocation);
}

public abandon() : void {
  this.tourExecution!.status = TourExecutionStatus.ABANDONED
  this.tourExecution!.finishedAt = new Date()  
  //http here 
}

/**
 * Haversine formula to calculate distance between two lat/lng points in meters
 */
private getDistanceMeters(lat1: number, lng1: number, lat2: number, lng2: number): number {
  const R = 6371000; // radius of Earth in meters
  const toRad = (x: number) => x * Math.PI / 180;

  const dLat = toRad(lat2 - lat1);
  const dLng = toRad(lng2 - lng1);

  const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
            Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) *
            Math.sin(dLng/2) * Math.sin(dLng/2);

  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  return R * c;
}

  private fetchTour(id: string): void {
    if (!this.savedLocation) return;

    this.http.get<Tour>(`http://localhost:8080/tour/${id}`).subscribe({
      next: (data) => {
        this.activatedTour = data;

        // --- CHECK IF TOUR EXECUTION EXISTS ---
        // TODO: fetch from backend, create if doesn't exist
        // this.tourExecution = fetchTourExecutionOrCreate(id);

        // For now, create dummy execution for testing
        this.tourExecution = {
          tourId: id,
          touristId: '',
          status: TourExecutionStatus.ACTIVE,
          startedAt: new Date(),
          lastActivityAt: new Date(),
          finishedAt: null,
          completedKeyPoints: [],
          currentPosition: {
            lat: this.savedLocation!.latitude,
            lng: this.savedLocation!.longitude
          }
        };

        // Draw uncompleted keypoints on the map
        this.drawKeypoints();
      },
      error: (err) => console.error('Error fetching tour:', err)
    });
  }

  private drawKeypoints(): void {
    if (!this.activatedTour || !this.tourExecution) return;

    // Clear old keypoint markers
    this.keypointMarkers.forEach(marker => this.map.removeLayer(marker));
    this.keypointMarkers = [];

    const completedIds = new Set(this.tourExecution.completedKeyPoints.map(kp => kp.keyPointId));

    // Assuming activatedTour has keypoints array like { id: string, latitude: number, longitude: number }
    this.activatedTour.keypoints
      .filter(kp => !completedIds.has(kp!.id!))
      .forEach(kp => {
        const marker = L.marker([kp.latitude, kp.longitude], { icon: this.keypointIcon })
          .addTo(this.map)
        this.keypointMarkers.push(marker);
      });
  }

}
