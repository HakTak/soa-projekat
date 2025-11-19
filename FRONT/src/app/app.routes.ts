import { Routes } from '@angular/router';
import { HomeComponent } from '../home/home.component';
import { TourCreateComponent } from './TOURS/tour-create/tour-create.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomeComponent },
    { path: 'tours', component: TourCreateComponent }
];
