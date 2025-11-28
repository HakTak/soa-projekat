import { Component, HostListener, OnInit } from '@angular/core';
import { RouterLink, Router } from "@angular/router";
import { CommonModule } from '@angular/common';
import { AuthService } from '../infrastructure/auth.service';
import { ShoppingCartService } from '../shopping-cart/service/shopping-cart.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink, CommonModule],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {
  cartCount = 0;
  dropdownOpen = false;
  isLoggedIn = false;

  constructor(private router: Router, private authService: AuthService, private cartService: ShoppingCartService) { }

  ngOnInit(): void {
    this.authService.userState$.subscribe((state: boolean) => {
      this.isLoggedIn = state;
    });

    this.cartService.cartCount$.subscribe(count => {
      this.cartCount = count;
    })
  }

  onLogout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  toggleDropdown(event: Event) {
    event.preventDefault();
    event.stopPropagation();
    this.dropdownOpen = !this.dropdownOpen;
  }

  closeDropdown() {
    this.dropdownOpen = false;
  }

  @HostListener('document:click')
  onDocumentClick() {
    this.dropdownOpen = false;
  }

  onHomeClick() {
    this.router.navigate(['/']);
  }

  onBlogClick() {
    this.router.navigate(['/blog']);
  }

  // --- POSTOJEĆI GETTERI ---
  get isAdmin(): boolean {
    return this.authService.hasRole() == "ADMIN";
  }

  get isTourist(): boolean {
    return this.authService.hasRole() == "TOURIST";
  }

  // --- NOVI GETTER (DODAJ OVO) ---
  get isGuide(): boolean {
    // Proverava da li je uloga 'GUIDE' (ili 'AUTHOR', zavisi kako ti se zove u bazi)
    return this.authService.hasRole() == "GUIDE"; 
  }
}