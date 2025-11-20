import { Routes } from '@angular/router';
import { HomeComponent } from '../home/home.component';
import { TourCreateComponent } from '../tour/tour-create/tour-create.component';
import { BlogComponent } from './blog/blog.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomeComponent },
    { path: 'tours', component: TourCreateComponent },
    { path: 'blog', component: BlogComponent }
];
