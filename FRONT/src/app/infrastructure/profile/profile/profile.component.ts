import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ProfileService } from '../service/profile.service';
import { Profile, UpdateProfileRequest, ProfileStats } from '../model/profile';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule], // Import modules here
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  profile: Profile | null = null;
  stats: ProfileStats | null = null;
  profileForm: FormGroup;

  isEditing = false;
  isLoading = false;
  errorMessage = '';

  readonly MAX_FILE_SIZE = 2 * 1024 * 1024;

  constructor(
    private profileService: ProfileService,
    private fb: FormBuilder
  ) {
    // Initialize form with validation
    this.profileForm = this.fb.group({
      firstName: ['', Validators.required],
      lastName: ['', Validators.required],
      profilePicture: [''],
      biography: [''],
      motto: ['']
    });
  }

  ngOnInit(): void {
    this.loadProfile();
  }

  loadProfile(): void {
    this.isLoading = true;
    this.profileService.getProfile().subscribe({
      next: (data) => {
        this.profile = data;
        this.fillForm(data);
        if (data.userId) {
          this.loadStats(data.userId);
        } else {
          this.isLoading = false;
        }
      },
      error: (err) => {
        console.error('Failed to load profile', err);
        this.errorMessage = 'Could not load profile data.';
        this.isLoading = false;
      }
    });
  }

  fillForm(data: Profile): void {
    this.profileForm.patchValue({
      firstName: data.firstName,
      lastName: data.lastName,
      profilePicture: data.profilePicture,
      biography: data.biography,
      motto: data.motto
    });
  }

  toggleEdit(): void {
    this.isEditing = !this.isEditing;
    if (!this.isEditing && this.profile) {
      // If user cancelled, reset form to original values
      this.fillForm(this.profile);
    }
  }

  onSubmit(): void {
    if (this.profileForm.invalid) return;

    this.isLoading = true;

    // Prepare object strictly typed
    const req: UpdateProfileRequest = {
      firstName: this.profileForm.value.firstName,
      lastName: this.profileForm.value.lastName,
      profilePicture: this.profileForm.value.profilePicture,
      biography: this.profileForm.value.biography,
      motto: this.profileForm.value.motto
    };

    this.profileService.updateProfile(req).subscribe({
      next: (updatedProfile) => {
        this.profile = updatedProfile;
        this.isEditing = false;
        this.isLoading = false;
        this.errorMessage = ''; // Clear errors
      },
      error: (err) => {
        console.error('Update failed', err);
        this.errorMessage = 'Failed to update profile. Please try again.';
        this.isLoading = false;
      }
    });
  }

  loadStats(userId: string): void {
    this.profileService.getStats(userId).subscribe({
      next: (statsData) => {
        this.stats = statsData;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Failed to load stats', err);
        this.isLoading = false;
      }
    });
  }

  // 1. TRIGGER FILE INPUT
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;

    const file = input.files[0];

    // 2. VALIDATION
    if (!file.type.match(/image\/*/)) {
      this.errorMessage = 'Only image files are allowed!';
      return;
    }

    if (file.size > this.MAX_FILE_SIZE) {
      this.errorMessage = 'Image size must be less than 2MB.';
      return;
    }

    // 3. CONVERT TO BASE64
    const reader = new FileReader();
    reader.onload = (e: any) => {
      const base64String = e.target.result;

      // Update the form field
      this.profileForm.patchValue({
        profilePicture: base64String
      });

      // Mark as dirty so the UI knows we changed something
      this.profileForm.markAsDirty();
    };

    reader.readAsDataURL(file);
  }
}