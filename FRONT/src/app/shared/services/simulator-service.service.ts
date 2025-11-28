import { Injectable } from '@angular/core';
import { Location } from '../../models/location.model';

@Injectable({
  providedIn: 'root'
})
export class SimulatorService {
  private storageKey = 'simulated_location';

  constructor() { }

  // Čuva lokaciju u LocalStorage
  setLocation(location: Location): void {
    localStorage.setItem(this.storageKey, JSON.stringify(location));
  }

  // Vraća sačuvanu lokaciju ili null ako ne postoji
  getLocation(): Location | null {
    const data = localStorage.getItem(this.storageKey);
    if (data) {
      try {
        return JSON.parse(data) as Location;
      } catch (e) {
        console.error('Error parsing location data', e);
        return null;
      }
    }
    return null;
  }

  // Briše lokaciju (opciono, ako zatreba reset)
  clearLocation(): void {
    localStorage.removeItem(this.storageKey);
  }
}