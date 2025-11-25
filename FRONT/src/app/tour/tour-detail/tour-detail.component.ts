import { Component, OnInit } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { Review } from '../model/review';
import { Tour } from '../model/tour';
import { Keypoint } from '../model/keypoint';

interface ReviewRequest {
  tour_id: string;
  user_id: string;
  user_name: string;
  rating: number;
  comment: string;
  created_at: string;
  visited_at: string;
  image_url: string | null;
}

@Component({
  selector: 'app-tour-detail',
  standalone: true,
  imports: [CommonModule, HttpClientModule, FormsModule],
  templateUrl: './tour-detail.component.html',
  styleUrls: ['./tour-detail.component.css']
})
export class TourDetailComponent implements OnInit {
  tour: Tour | null = null;
  reviews: Review[] = [];
  tourImage: string | null = null;
  
  // Modal Kontrola
  isModalOpen = false;

  // Notification Popup Kontrola
  showNotification = false;
  notificationMessage = '';
  notificationType: 'success' | 'error' = 'success';

  // Podaci za novi review
  newReview: ReviewRequest = {
    tour_id: '',
    user_id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
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
        this.setTourImage();
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
    this.router.navigate(['/tour-list']); // Proveri da li je ovo tvoja ispravna ruta za listu
  }

  // --- DELETE FUNKCIONALNOST ---

  deleteTour(): void {
    if (!this.tour || !this.tour.id) return;

    // Potvrda od korisnika pre brisanja
    const confirmDelete = confirm(`Are you sure you want to delete tour "${this.tour.title}"? This action cannot be undone.`);
    
    if (confirmDelete) {
      this.http.delete(`http://localhost:8083/tour/${this.tour.id}`).subscribe({
        next: () => {
          // Prikazi uspeh
          this.showToast('Tour successfully deleted!', 'success');
          
          // Vrati nazad na listu posle 2 sekunde (da stigne da vidi poruku)
          setTimeout(() => {
            this.router.navigate(['/tour-list']);
          }, 2000);
        },
        error: (err) => {
          console.error('Delete error:', err);
          this.showToast('Failed to delete tour. Please try again.', 'error');
        }
      });
    }
  }

  // --- TOAST NOTIFICATION HELPER ---
  showToast(message: string, type: 'success' | 'error'): void {
    this.notificationMessage = message;
    this.notificationType = type;
    this.showNotification = true;

    // Sakrij automatski posle 3 sekunde
    setTimeout(() => {
      this.showNotification = false;
    }, 3000);
  }

  // --- LOGIKA ZA MODAL ---

  openReviewModal(): void {
    if (!this.tour) return;
    
    this.newReview = {
      tour_id: this.tour.id ?? '',
      user_id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
      user_name: 'Petar',
      rating: 5,
      comment: '',
      created_at: '',
      visited_at: new Date().toISOString().split('T')[0],
      image_url: null
    };
    
    this.isModalOpen = true;
  }

  closeReviewModal(): void {
    this.isModalOpen = false;
  }

  submitReview(): void {
    const payload = {
      ...this.newReview,
      created_at: new Date().toISOString(),
      visited_at: new Date(this.newReview.visited_at).toISOString(),
      rating: Number(this.newReview.rating)
    };

    this.http.post('http://localhost:8083/review/create', payload).subscribe({
      next: (res) => {
        this.showToast('Review successfully saved!', 'success'); // Koristimo novi toast umesto alert-a
        this.closeReviewModal();
        if (this.tour) {
          this.fetchReviews(this.tour.id ?? '');
        }
      },
      error: (err) => {
        console.error('Error creating review:', err);
        this.showToast('Failed to save review.', 'error');
      }
    });
  }
}