import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavbarComponent } from "./navbar/navbar.component";
import { HomepageComponent } from './homepage/homepage.component';
import { BlogComponent } from './blog/blog.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    NavbarComponent,
    HomepageComponent,
    BlogComponent

  ],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']   // ✔️ FIXED
})
export class AppComponent {
  title = 'FRONT';
}
