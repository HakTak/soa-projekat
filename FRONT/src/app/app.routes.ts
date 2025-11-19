import { Routes } from '@angular/router';
import { HomepageComponent } from './homepage/homepage.component';
import { BlogComponent } from './blog/blog.component';

export const routes: Routes = [
    { path: '', component: HomepageComponent },
    { path:'blog', component: BlogComponent }
];
