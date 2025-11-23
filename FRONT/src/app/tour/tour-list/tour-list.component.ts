import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Tour } from '../model/tour'; // Prilagodi putanju ako je drugacija
import { Keypoint } from '../model/keypoint'; // Prilagodi putanju ako je drugacija

@Component({
  selector: 'app-tour-list',
  standalone: true,
  imports: [CommonModule, HttpClientModule],
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  tours: Tour[] = [];

  constructor(private http: HttpClient, private router: Router) {}

  ngOnInit(): void {
    this.fetchTours();
  }

  fetchTours(): void {
    this.http.get<Tour[]>('http://localhost:8083/tours').subscribe({
      next: (data) => {
        this.tours = data;
      },
      error: (err) => {
        console.error('Error fetching tours:', err);
      }
    });
  }

  // Pomocna metoda za navigaciju
  onCardClick(tourId: string | undefined): void {
    if (tourId) {
      // Ovde menjas putanju gde zelis da te odvede klik (npr. /tours/123)
      this.router.navigate(['/tours', tourId]);
    }
  }

  // Dobija ime prve kljucne tacke
  getFirstKeypointName(tour: Tour): string {
    return tour.keypoints && tour.keypoints.length > 0 
      ? tour.keypoints[0].title 
      : 'No start point';
  }

  // Dobija ime poslednje kljucne tacke
  getLastKeypointName(tour: Tour): string {
    return tour.keypoints && tour.keypoints.length > 0 
      ? tour.keypoints[tour.keypoints.length - 1].title 
      : 'No end point';
  }
}