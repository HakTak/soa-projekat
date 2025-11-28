import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Tour } from '../model/tour';
import { Keypoint } from '../model/keypoint'; // Ako ti treba
import { ShoppingCartService } from '../../shopping-cart/service/shopping-cart.service';
import { AuthService } from '../../infrastructure/auth.service';
import { TourStatus } from '../model/tour';

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
  public tours: Tour[] = [];
  public tourStatus = TourStatus;
  
  // Mapa koja cuva ocene: Kljuc je ID ture, Vrednost je prosek (number)
  tourRatings: { [tourId: string]: number } = {};

  constructor(private http: HttpClient, private router: Router, private cartService: ShoppingCartService, private authService: AuthService) { }

  ngOnInit(): void {
    this.fetchTours();
    this.tours = this.tours.filter(tour => tour.status === TourStatus.PUBLISHED);

  }

  fetchTours(): void {
    const headers = this.authService.getAuthHeaders();
    const roleLogged = this.authService.hasRole();
    const userId = this.authService.getMyId();
    const userName = this.authService.getMyUsername();
    console.log('Is tourist logged in?', userName);

    this.http.get<any>('http://localhost:8080/tours', { headers: headers }).subscribe({
      next: (response) => {
        // ISPRAVKA: Vadimo niz iz polja 'tours'
        // Ako response.tours ne postoji, stavljamo prazan niz [] da ne pukne app
        const rawTours = response.tours || [];


        if (roleLogged === 'TOURIST') { // dodati filter samo publiched
          this.tours = rawTours
          console.log('Fetched tours for TOURIST:', rawTours);
          return;
        }

        if (roleLogged === 'ADMIN') {
          this.tours = rawTours; // Admin vidi sve ture
          console.log('Fetched tours for ADMIN:', rawTours);
          return;
        }

        if (roleLogged === 'AUTHOR') { //Samo ture od autora
          this.tours = rawTours.filter((tour: Tour) => (tour.userName === userName));
          console.log('Fetched tours for AUTHOR:', rawTours);
          return;
        }
        // 1. Sada filtriramo taj niz (ne response objekat, nego niz unutar njega)
       this.tours = rawTours.filter((tour: Tour) => tour.status === TourStatus.DRAFT);
       console.log('Fetched tours:', rawTours);

        // 2. Dovlačimo ocene
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
    const headers = this.authService.getAuthHeaders();
    this.http.get<Review[]>(`http://localhost:8080/review/tour/${tourId}`, { headers: headers }).subscribe({
      next: (reviews) => {
        console.log(`Fetched reviews for tour ${tourId}:`, reviews);
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

  addToCart(event: Event, tourId?: string): void {
    event.stopPropagation();

    const userId = this.authService.getMyId();
    if (!userId) {
      alert('Please log in to shop.');
      return;
    }

    this.cartService.addItem(userId, tourId).subscribe({
      next: () => {
        // Optional: Show a toast notification
        alert('Added to cart!');
      },
      error: (err) => {
        alert(err.error?.message || 'Failed to add item. It might already be in your cart.');
      }
    });
  }
}