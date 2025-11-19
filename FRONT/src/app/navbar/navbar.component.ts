import { Component } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { Route } from '@angular/router';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})

export class NavbarComponent {
toggleNavbar() {
throw new Error('Method not implemented.');
}
isOpen: any;

 constructor(private router: Router) {}

onHomeClick() {
    console.log("home clicked!");
    this.router.navigate(['/']);
  }

onBlogClick() {
    console.log("blog clicked!");
    this.router.navigate(['/blog']);
  }
}
