import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Tour } from '../model/tour'; 
 import { Keypoint } from '../model/keypoint'; // Ako ti treba

// Definisemo mali interfejs za Review ovde (ili ga importuj ako ga imas u models folderu)
interface Review {
  id?: string;
  rating: number;
  comment?: string;
  tourId: string;
}

@Component({
  selector: 'app-tour-list',
  standalone: true,
  imports: [CommonModule, HttpClientModule],
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  tours: Tour[] = [];
  
  // Mapa koja cuva ocene: Kljuc je ID ture, Vrednost je prosek (number)
  tourRatings: { [tourId: string]: number } = {};

  constructor(private http: HttpClient, private router: Router) {}

  ngOnInit(): void {
    this.fetchTours();
  }

  fetchTours(): void {
    // Koristimo relativnu putanju jer imamo Proxy (ili CORS setup)
    this.http.get<Tour[]>('http://localhost:8083/tours').subscribe({
      next: (data) => {
        this.tours = data;
        
        // Cim stignu ture, za svaku od njih pokreni dovlacenje ocena
        this.tours.forEach(tour => {
          if (tour.id) {
            this.fetchAvgRating(tour.id);
          }
        });
      },
      error: (err) => {
        console.error('Error fetching tours:', err);
      }
    });
  }

  fetchAvgRating(tourId: string): void {
    this.http.get<Review[]>(`http://localhost:8083/review/tour/${tourId}`).subscribe({
      next: (reviews) => {
        if (reviews.length > 0) {
          // Izracunaj sumu
          const sum = reviews.reduce((acc, review) => acc + review.rating, 0);
          // Izracunaj prosek
          const avg = sum / reviews.length;
          // Sacuvaj u mapu pod ID-jem ture
          this.tourRatings[tourId] = avg;
        } else {
          // Ako nema recenzija, stavimo 0
          this.tourRatings[tourId] = 0;
        }
      },
      error: (err) => {
        console.error(`Error fetching reviews for tour ${tourId}:`, err);
        this.tourRatings[tourId] = 0; // Fallback na 0 ako pukne request
      }
    });
  }

  onCardClick(tourId: string | undefined): void {
    if (tourId) {
      this.router.navigate(['/tours', tourId]);
    }
  }

  getFirstKeypointName(tour: Tour): string {
    return tour.keypoints && tour.keypoints.length > 0 
      ? tour.keypoints[0].title 
      : 'No start point';
  }

  getLastKeypointName(tour: Tour): string {
    return tour.keypoints && tour.keypoints.length > 0 
      ? tour.keypoints[tour.keypoints.length - 1].title 
      : 'No end point';
  }
}