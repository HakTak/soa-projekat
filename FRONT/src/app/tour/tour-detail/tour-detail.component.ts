import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient, HttpClientModule, HttpHeaders } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { Review } from '../model/review';
import { Tour, TourStatus } from '../model/tour';
import { Keypoint } from '../model/keypoint';
import { User } from '../../models/user.model';
import { AuthService } from '../../infrastructure/auth.service';
import { ShoppingCartService } from '../../shopping-cart/service/shopping-cart.service';

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
  tourStatus = TourStatus;
  selectedFileName: string | null = null;
  
  // Status kupovine
  isPurchased: boolean = false; 

  // Modal & Notification
  isModalOpen = false;
  showNotification = false;
  notificationMessage = '';
  notificationType: 'success' | 'error' = 'success';

  // Novi review
  newReview: ReviewRequest = {
    tour_id: '',
    user_id: '',
    user_name: '',
    rating: 5,
    comment: '',
    created_at: '',
    visited_at: '',
    image_url: null
  };

  public user: User | null = this.loadCurrentUser();

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    private router: Router,
    private authService: AuthService,
    private shoppingCartService: ShoppingCartService
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

  loadCurrentUser() {
   const token = localStorage.getItem('jwt');
    if (!token) return null;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return {
        id: payload.id,
        username: payload.username,
        email: payload.email,
        role: payload.role,
        blocked: false
      } as User;
    } catch (err) {
      console.error("Invalid JWT:", err);
      return null;
    }
  }

  fetchTour(id: string): void {
    const headers = this.authService.getAuthHeaders();
    this.http.get<any>(`http://localhost:8080/tour/${id}`, { headers: headers }).subscribe({
      next: (data) => {
        console.log("Fetched tour data:", data);
        this.tour = data.tour;

        // --- FIX ZA DUPLIKATE ---
        // Posto baza vraca duplirane tacke, ovde ih filtriramo da ostanu samo unikatne (po ID-u ili Naslovu)
        if (this.tour && this.tour.keypoints) {
            this.tour.keypoints = this.tour.keypoints.filter((kp, index, self) =>
                index === self.findIndex((t) => (
                    t.id === kp.id || t.title === kp.title
                ))
            );
        }
        // ------------------------

        this.setTourImage();

        if (this.user && this.user.role === 'TOURIST') {
          this.checkPurchaseStatus();
        }
      },
      error: (err) => console.error('Error fetching tour:', err)
    });
  }

  checkPurchaseStatus(): void {
    if (!this.user || !this.user.id || !this.tour) return;

    const headers = this.authService.getAuthHeaders();
    this.http.get<any>(`http://localhost:8080/api/shopping-cart/orders/${this.user.id}`, { headers: headers })
      .subscribe({
        next: (response) => {
          const purchasedTokens = response.tokens || response || [];
          const match = purchasedTokens.find((token: any) => 
            (token.tourId === this.tour?.id) || (token.tour_id === this.tour?.id)
          );
          this.isPurchased = !!match;
        },
        error: (err) => {
          console.error("Error checking purchase status:", err);
          this.isPurchased = false;
        }
      });
  }

  fetchReviews(tourId: string): void {
    const headers = this.authService.getAuthHeaders();
    this.http.get<any>(`http://localhost:8080/review/tour/${tourId}`, { headers }).subscribe({
      next: (data) => {
        const reviewsArray = data.reviews || []; 
        this.reviews = reviewsArray.map((r: any) => {
          return {
            id: r.id,
            rating: r.rating,
            comment: r.comment,
            username: r.userName,
            createdAt: r.createdAt,
            imageUrl: r.imageUrl
          };
        });
      },
      error: (err) => console.error('Error loading reviews:', err)
    });
  }

  // --- PUBLISH / ARCHIVE (VRACENO NA STARO) ---
  
  publishTour(): void { 
    if (this.tour == null) return;
    
    this.tour.status = this.tourStatus.PUBLISHED;
    this.tour.publisedAt = new Date();
    this.tour.archivedAt = null; 
    
    // Saljemo ceo this.tour objekat, bez posebnog payload-a
    this.http.patch<any>('http://localhost:8080/tour/update', this.tour, {headers: this.getAuthHeaders()}).subscribe({
      next: () => {
        this.showToast('Tour successfully published!', 'success');
      },
      error: (err) => {
        console.error('Publish error:', err);
        this.showToast('Failed to publish tour.', 'error');
        if (this.tour) this.tour.status = this.tourStatus.DRAFT; 
      }
    })
  }

  archiveTour(): void {
    if (this.tour == null) return;
    
    this.tour.status = this.tourStatus.ARCHIVED;
    this.tour.archivedAt = new Date();
    this.tour.publisedAt = null;

    // Saljemo ceo this.tour objekat
    this.http.patch<any>('http://localhost:8080/tour/update', this.tour, {headers: this.getAuthHeaders()}).subscribe({
      next: () => {
        this.showToast('Tour successfully archived!', 'success');
      },
      error: (err) => {
        console.error('Archive error:', err);
        this.showToast('Failed to archive tour.', 'error');
        if (this.tour) this.tour.status = this.tourStatus.PUBLISHED;
      }
    })
  }

  // --- OSTALE POMOCNE FUNKCIJE ---

  onFileSelected(event: any): void {
    const file: File = event.target.files[0];
    if (file) {
      this.selectedFileName = file.name;
      const reader = new FileReader();
      reader.onload = (e: any) => {
        this.newReview.image_url = e.target.result; 
      };
      reader.readAsDataURL(file);
    }
  }

  removeImage(): void {
    this.newReview.image_url = null;
    this.selectedFileName = null;
    const fileInput = document.getElementById('fileInput') as HTMLInputElement;
    if(fileInput) fileInput.value = '';
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

  deleteTour(): void {
    if (!this.tour || !this.tour.id) return;
    const confirmDelete = confirm(`Are you sure you want to delete tour "${this.tour.title}"?`);
    
    if (confirmDelete) {
      this.http.delete(`http://localhost:8080/tour/${this.tour.id}`, {headers: this.getAuthHeaders()}).subscribe({
        next: () => {
          this.showToast('Tour successfully deleted!', 'success');
          setTimeout(() => {
            this.router.navigate(['/tour-list']);
          }, 2000);
        },
        error: (err) => {
          console.error('Delete error:', err);
          this.showToast('Failed to delete tour.', 'error');
        }
      });
    }
  }

  openReviewModal(): void {
    if (!this.tour) return;
    const userId = this.user?.id ?? '';
    const userName = this.user?.username ?? 'Anonymous';
    this.newReview = {
      tour_id: this.tour.id ?? '',
      user_id: userId,
      user_name: userName,
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

    this.http.post('http://localhost:8080/review/create', payload, {headers: this.getAuthHeaders()}).subscribe({
      next: (res) => {
        this.showToast('Review saved!', 'success');
        this.closeReviewModal();
        if (this.tour) this.fetchReviews(this.tour.id ?? '');
      },
      error: (err) => {
        console.error('Error creating review:', err);
        this.showToast('Failed to save review.', 'error');
      }
    });
  }

  showToast(message: string, type: 'success' | 'error'): void {
    this.notificationMessage = message;
    this.notificationType = type;
    this.showNotification = true;
    setTimeout(() => {
      this.showNotification = false;
    }, 3000);
  }

  getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt');;
    return new HttpHeaders({ 'Authorization': `Bearer ${token}` });
  }

    addToCart(event: Event, tourId?: string): void {
    event.stopPropagation();

    const userId = this.authService.getMyId();
    if (!userId) {
      alert('Please log in to shop.');
      return;
    }

    this.shoppingCartService.addItem(userId, tourId).subscribe({
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