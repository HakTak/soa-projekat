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
    this.tours = [
    {
      id: "t1",
      title: "Alps Mountain Adventure",
      description: "A thrilling multi-day guided tour across the Alpine passes.",
      difficulty: "Hard",
      tags: "mountains, hiking, adventure",
      status: "published",
      publisedAt: new Date(),
      archivedAt: null,
      price: 299.99,
      keypoints: [
        {
          id: "k1",
          tourId: "t1",
          title: "Base Camp",
          latitude: 46.8182,
          longitude: 8.2275,
          description: "Starting point at the scenic Swiss Alps base camp.",
          imageUrl: "https://picsum.photos/300/200?alps1"
        },
        {
          id: "k2",
          tourId: "t1",
          title: "Glacier Ridge",
          latitude: 46.9121,
          longitude: 7.9980,
          description: "A stunning viewpoint overlooking an ancient glacier.",
          imageUrl: "https://picsum.photos/300/200?alps2"
        }
      ]
    },
    {
      id: "t2",
      title: "Old Town Cultural Walk",
      description: "A relaxed guided walk through medieval European streets.",
      difficulty: "Easy",
      tags: "city, history, culture",
      status: "published",
      publisedAt: new Date(),
      archivedAt: null,
      price: 49.99,
      keypoints: [
        {
          id: "k3",
          tourId: "t2",
          title: "Town Square",
          latitude: 48.8566,
          longitude: 2.3522,
          description: "The historic square surrounded by gothic architecture.",
          imageUrl: "https://picsum.photos/300/200?city1"
        },
        {
          id: "k4",
          tourId: "t2",
          title: "Ancient Library",
          latitude: 48.8575,
          longitude: 2.3500,
          description: "Home to rare manuscripts from the 12th century.",
          imageUrl: "https://picsum.photos/300/200?city2"
        }
      ]
    },
    {
      id: "t3",
      title: "Island Kayak Escape",
      description: "Paddle between remote islands and secret beaches.",
      difficulty: "Medium",
      tags: "water, beach, kayaking",
      status: "draft",
      publisedAt: null,
      archivedAt: null,
      price: 0,
      keypoints: [
        {
          id: "k5",
          tourId: "t3",
          title: "Lagoon Start Point",
          latitude: 36.3932,
          longitude: 25.4615,
          description: "Turquoise shallow waters perfect for beginners.",
          imageUrl: "https://picsum.photos/300/200?sea1"
        },
        {
          id: "k6",
          tourId: "t3",
          title: "Hidden Beach",
          latitude: 36.3910,
          longitude: 25.4650,
          description: "A quiet beach only accessible by kayak.",
          imageUrl: "https://picsum.photos/300/200?sea2"
        }
      ]
    },
    {
      id: "t4",
      title: "Rainforest Wildlife Trek",
      description: "A guided adventure through lush rainforest trails.",
      difficulty: "Hard",
      tags: "forest, wildlife, nature",
      status: "published",
      publisedAt: new Date(),
      archivedAt: null,
      price: 159.50,
      keypoints: [
        {
          id: "k7",
          tourId: "t4",
          title: "Waterfall Point",
          latitude: -3.4653,
          longitude: -62.2159,
          description: "A massive waterfall surrounded by dense vegetation.",
          imageUrl: "https://picsum.photos/300/200?forest1"
        },
        {
          id: "k8",
          tourId: "t4",
          title: "Wildlife Lookout",
          latitude: -3.4700,
          longitude: -62.2100,
          description: "Perfect spot for observing tropical bird species.",
          imageUrl: "https://picsum.photos/300/200?forest2"
        }
      ]
    }
  ];
  this.tours = this.tours.filter(tour => tour.status === 'published');

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