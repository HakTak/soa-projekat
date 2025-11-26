import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

import { User } from '../../models/user.model';
// Pazi: Importujemo NOVI servis
import { AuthFeatureService } from '../auth-service.service';

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './user-list.component.html',
  styleUrl: './user-list.component.css'
})
export class UserListComponent implements OnInit {

  users: User[] = [];

  // Injectujemo AuthFeatureService
  constructor(private userService: AuthFeatureService, private router: Router) { }

  ngOnInit(): void {
    this.loadUsers();
  }

  loadUsers() {
    this.userService.getUsers().subscribe({
      next: (response) => {
        this.users = response.users || [];
      },
      error: (err) => {
        console.error(err);
        alert('Greska ili nemate prava pristupa.');
        this.router.navigate(['/home']);
      }
    });
  }

  onBlock(user: User) {
    if (!confirm(`Blokirati korisnika ${user.username}?`)) return;

    this.userService.blockUser(user.id).subscribe({
      next: () => {
        alert('Blokiran!');
        user.blocked = true; // Azuriramo prikaz lokalno
      },
      error: (err) => alert('Greska pri blokiranju.')
    });
  }
}