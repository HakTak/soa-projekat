import { Routes } from '@angular/router';
import { HomepageComponent } from './homepage/homepage.component';
import { TourCreateComponent } from './tour/tour-create/tour-create.component';
import { BlogComponent } from './blog/blog-create/blog.component';
import { BlogsComponent } from './blog/blogs/blogs.component';
import { TourListComponent } from './tour/tour-list/tour-list.component';
import { TourDetailComponent } from './tour/tour-detail/tour-detail.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomepageComponent },
    { path: 'tours', component: TourCreateComponent },
    { path: 'blog/create', component: BlogComponent },
    { path: 'blogs', component: BlogsComponent },
    {path: 'tour-list', component: TourListComponent},
    {path: 'tours/:id', component: TourDetailComponent}
];
