import { Component, OnInit } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms'; // <--- OBAVEZNO DODATI
import { Review } from '../model/review';
import { Tour } from '../model/tour';
import { Keypoint } from '../model/keypoint';

// Model za formu
interface ReviewRequest {
  tour_id: string;
  user_id: string;
  user_name: string;
  rating: number;
  comment: string;
  created_at: string;
  visited_at: string; // Ovo ce biti string iz date inputa
  image_url: string | null;
}

@Component({
  selector: 'app-tour-detail',
  standalone: true,
  imports: [CommonModule, HttpClientModule, FormsModule], // <--- FormsModule dodat
  templateUrl: './tour-detail.component.html',
  styleUrls: ['./tour-detail.component.css']
})
export class TourDetailComponent implements OnInit {
  tour: Tour | null = null;
  reviews: Review[] = [];
  tourImage: string | null = null;
  
  // Kontrola modala
  isModalOpen = false;

  // Podaci za novi review
  newReview: ReviewRequest = {
    tour_id: '',
    user_id: '3fa85f64-5717-4562-b3fc-2c963f66afa6', // Hardcodovan ID (simulacija ulogovanog usera)
    user_name: 'Petar',
    rating: 5,
    comment: '',
    created_at: '',
    visited_at: '',
    image_url: null
  };

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const id = params.get('id');
      if (id) {
        this.fetchTour(id);
        this.fetchReviews(id);
      }
    });
  }

  fetchTour(id: string): void {
    this.http.get<Tour>(`http://localhost:8083/tour/${id}`).subscribe({
      next: (data) => {
        this.tour = data;
        //this.setTourImage();
      },
      error: (err) => console.error('Error fetching tour:', err)
    });
  }

  fetchReviews(tourId: string): void {
    this.http.get<any[]>(`http://localhost:8083/review/tour/${tourId}`).subscribe({
  next: (data) => {
    this.reviews = data.map(r => ({
      id: r.id,
      rating: r.rating,
      comment: r.comment,
      tourId: r.tour_id,
      userId: r.user_id,
      username: r.user_name,
      createdAt: r.created_at ? new Date(r.created_at) : undefined
    }));
    console.log('Mapped reviews:', this.reviews);
  },
  error: (err) => console.error('Error fetching reviews:', err)
});
  }

  setTourImage(): void {
    if (this.tour && this.tour.keypoints) {
      const validPoint = this.tour.keypoints.find(kp => kp.imageUrl && kp.imageUrl.trim() !== '');
      this.tourImage = validPoint ? validPoint.imageUrl! : null;
    }
  }

  goBack(): void {
    this.router.navigate(['/tour-list']);
  }

  // --- LOGIKA ZA MODAL ---

  openReviewModal(): void {
    if (!this.tour) return;
    
    // Resetuj formu pre otvaranja
    this.newReview = {
      tour_id: this.tour.id??'',
      user_id: '3fa85f64-5717-4562-b3fc-2c963f66afa6', // Hardcodovan user ID
      user_name: 'Petar',
      rating: 5,
      comment: '',
      created_at: '',
      visited_at: new Date().toISOString().split('T')[0], // Danasnji datum kao default
      image_url: null
    };
    
    this.isModalOpen = true;
  }

  closeReviewModal(): void {
    this.isModalOpen = false;
  }

  submitReview(): void {
    // 1. Pripremi podatke (formatiraj datume u ISO format kako backend ocekuje)
    const payload = {
      ...this.newReview,
      created_at: new Date().toISOString(), // Trenutno vreme
      visited_at: new Date(this.newReview.visited_at).toISOString(), // Datum posete + 00:00:00
      rating: Number(this.newReview.rating) // Osiguraj da je rating broj
    };

    // 2. Posalji POST zahtev
    console.log('Submitting review payload:', payload);
    this.http.post('http://localhost:8083/review/create', payload).subscribe({
      next: (res) => {
        alert('Review successfully saved!'); // Notifikacija
        this.closeReviewModal(); // Zatvori modal
        if (this.tour) {
          this.fetchReviews(this.tour.id ?? ''); // Osvezi listu komentara
        }
      },
      error: (err) => {
        console.error('Error creating review:', err);
        alert('Failed to save review. Check console.');
      }
    });
  }
}