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
    console.log('Is tourist logged in?', roleLogged);

    this.http.get<any>('http://localhost:8080/tours', { headers: headers }).subscribe({
      next: (response) => {
        // ISPRAVKA: Vadimo niz iz polja 'tours'
        // Ako response.tours ne postoji, stavljamo prazan niz [] da ne pukne app
        const rawTours = response.tours || [];


        if (roleLogged === 'TOURIST') { // dodati filter samo publiched
          this.tours = rawTours.filter((tour: Tour) => tour.status === TourStatus.PUBLISHED);
          console.log('Fetched tours for TOURIST:', rawTours);
          
        }

        if (roleLogged === 'ADMIN') {
          this.tours = rawTours; // Admin vidi sve ture
          console.log('Fetched tours for ADMIN:', rawTours);
          
        }

        if (roleLogged === 'GUIDE') { //Samo ture od autora
          this.tours = rawTours ;

          console.log('Fetched tours for AUTHOR:', rawTours);
          
        }
        // 1. Sada filtriramo taj niz (ne response objekat, nego niz unutar njega)
      // this.tours = rawTours.filter((tour: Tour) => tour.status === TourStatus.DRAFT);
       console.log('Fetched tours:', rawTours);

        // 2. Dovlačimo ocene
        console.log("ALOOOO 111.");
        this.tours.forEach(tour => {
          if (tour.id) {
            console.log("ALOOO 222:", tour.id);
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
      
      // 1. Koristimo <any> jer odgovor nije cist niz Review[], vec objekat
      this.http.get<any>(`http://localhost:8080/review/tour/${tourId}`, { headers: headers }).subscribe({
        next: (response) => {
          
          // 2. Izvlacimo niz iz polja 'reviews'
          // Ako polje ne postoji, koristimo prazan niz da ne pukne kod
          const reviewsList = response.reviews || []; 

          console.log(`Reviews za ${tourId}:`, reviewsList);

          if (reviewsList.length > 0) {
            // 3. Koristimo 'reviewsList' za racunanje, a ne 'response'
            const sum = reviewsList.reduce((acc: number, review: Review) => acc + review.rating, 0);
            const avg = sum / reviewsList.length;
            
            this.tourRatings[tourId] = avg;
          } else {
            this.tourRatings[tourId] = 0;
          }
        },
        error: (err) => {
          console.error(`Error fetching reviews for tour ${tourId}:`, err);
          this.tourRatings[tourId] = 0;
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