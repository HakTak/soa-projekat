import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router'; // <--- Router ide ovde
// Import iz foldera iznad (infrastructure)
import { AuthService, LoginRequest } from '../auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {

  loginData: LoginRequest = {
    username: '',
    password: ''
  };

  // Ubacujemo Router u konstruktor komponente
  constructor(private authService: AuthService, private router: Router) { }

  onLogin() {
    this.authService.login(this.loginData).subscribe({
      next: (res) => {
        console.log('Login uspesan');
        // Komponenta radi redirekciju
        this.router.navigate(['/home']);
      },
      error: (err) => {
        console.error(err);
        alert('Pogresni podaci ili greska na serveru');
      }
    });
  }
}