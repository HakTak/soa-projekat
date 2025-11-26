import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
// Import modela i servisa
import { ProfileFollowService } from '../profile-follow.service';
import { FullProfileDisplay } from '../../models/profile-follow.model';
import { AuthService } from '../../infrastructure/auth.service'; // Za proveru logina

@Component({
  selector: 'app-profile-follow-page',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './profile-follow-page.component.html',
  styleUrl: './profile-follow-page.component.css'
})
export class ProfileFollowPageComponent implements OnInit {

  userProfile: FullProfileDisplay | null = null;
  isLoading = true;
  isOwnProfile = false;

  constructor(
    private route: ActivatedRoute,
    private pfService: ProfileFollowService,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const userId = params.get('id');
      if (userId) {
        this.loadData(userId);
      }
    });
  }

  loadData(targetId: string) {
    this.isLoading = true;

    // Provera da li sam to ja (dekodiranjem tokena ili proverom ID-ja ako ga imamo u localStorage)
    // Ovde cemo jednostavnije: ako sam ulogovan, pitam bekend
    const isLoggedIn = this.authService.isLoggedIn();
    const myId = this.authService.getMyId();

    // 1. Dovuci Profil
    this.pfService.getProfile(targetId).subscribe({
      next: (p) => {
        this.userProfile = { ...p };

        // 2. Dovuci Statistiku
        this.pfService.getStats(targetId).subscribe(s => {
          if (this.userProfile) this.userProfile.stats = s;
        });

        // 3. Proveri Follow status (samo ako sam ulogovan)
        if (isLoggedIn) {
          // Provera da li je ovo moj profil bi se idealno radila poredjenjem ID-ja iz tokena
          // Ali za sad cemo probati isFollowing, ako backend vrati gresku ili false, znamo stanje
          this.pfService.isFollowing(myId, targetId).subscribe({
            next: (res) => {
              if (this.userProfile) this.userProfile.isFollowing = res.isFollowing;
              console.log("EVO: " + res.isFollowing)
            },
            error: () => {
              // Ako pukne, pretpostavimo false
              if (this.userProfile) this.userProfile.isFollowing = false;
            }
          });
        }
        this.isLoading = false;
      },
      error: (err) => {
        console.error(err);
        this.isLoading = false;
      }
    });
  }

  toggleFollow() {
    if (!this.userProfile) return;

    if (this.userProfile.isFollowing) {
      this.pfService.unfollow(this.userProfile.userId).subscribe(() => {
        this.userProfile!.isFollowing = false;
        if (this.userProfile!.stats) this.userProfile!.stats!.followersCount--;
      });
    } else {
      this.pfService.follow(this.userProfile.userId).subscribe(() => {
        this.userProfile!.isFollowing = true;
        if (this.userProfile!.stats) this.userProfile!.stats!.followersCount++;
      });
    }
  }
}