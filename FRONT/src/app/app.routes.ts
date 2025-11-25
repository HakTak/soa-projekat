import { Routes } from '@angular/router';
import { HomepageComponent } from './homepage/homepage.component';
import { TourCreateComponent } from './tour/tour-create/tour-create.component';
import { BlogComponent } from './blog/blog-create/blog.component';
import { BlogsComponent } from './blog/blogs/blogs.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomepageComponent },
    { path: 'tours', component: TourCreateComponent },
    { path: 'blog/create', component: BlogComponent },
    { path: 'blogs', component: BlogsComponent }
];
