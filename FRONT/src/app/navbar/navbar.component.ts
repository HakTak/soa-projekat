import { Component, HostListener } from '@angular/core';
import { RouterLink } from "@angular/router";
import { Router } from '@angular/router';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})

export class NavbarComponent {
  dropdownOpen = false;

  constructor(private router: Router) {}

  // Ova metoda sluzi za klik na samo "Tours" dugme
  toggleDropdown(event: Event) {
    event.preventDefault();
    event.stopPropagation();
    this.dropdownOpen = !this.dropdownOpen;
  }

  // Ova metoda sluzi za eksplicitno zatvaranje (npr. kada se klikne na Create Tour)
  closeDropdown() {
    this.dropdownOpen = false;
  }

  // Zatvaranje kada se klikne bilo gde drugde na stranici
  @HostListener('document:click')
  onDocumentClick() {
    this.dropdownOpen = false;
  }

  toggleNavbar() {
    // Ako se koristi za hamburger meni u buducnosti
    throw new Error('Method not implemented.');
  }

  isOpen: any;

  onHomeClick() {
    console.log("home clicked!");
    this.router.navigate(['/']);
  }

  onBlogClick() {
    console.log("blog clicked!");
    this.router.navigate(['/blog']);
  }
}