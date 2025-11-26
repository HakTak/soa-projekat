import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { ProfileFollowService } from '../profile-follow.service';
import { FullProfileDisplay } from '../../models/profile-follow.model';

@Component({
  selector: 'app-profile-recommended-page',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './profile-recommended-page.component.html',
  styleUrl: './profile-recommended-page.component.css'
})
export class ProfileRecommendedPageComponent implements OnInit {

  recommendations: FullProfileDisplay[] = [];
  isLoading = true;

  constructor(private pfService: ProfileFollowService) { }

  ngOnInit(): void {
    this.pfService.getEnrichedRecommendations().subscribe({
      next: (data) => {
        this.recommendations = data;
        this.isLoading = false;
      },
      error: (err) => {
        console.error("Failed to load recommendations", err);
        this.isLoading = false;
      }
    });
  }
}