import { Routes } from '@angular/router';
import { HomepageComponent } from './homepage/homepage.component';
import { TourCreateComponent } from './tour/tour-create/tour-create.component';
import { BlogComponent } from './blog/blog-create/blog.component';
import { BlogsComponent } from './blog/blogs/blogs.component';
import { TourListComponent } from './tour/tour-list/tour-list.component';
import { TourDetailComponent } from './tour/tour-detail/tour-detail.component';
import { LoginComponent } from './infrastructure/login/login.component'
import { RegisterComponent } from './infrastructure/register/register.component';
import { UserListComponent } from './auth/user-list/user-list.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomepageComponent },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'tours', component: TourCreateComponent },
    { path: 'blog/create', component: BlogComponent },
    { path: 'blogs', component: BlogsComponent },
    { path: 'tour-list', component: TourListComponent},
    { path: 'tours/:id', component: TourDetailComponent},
    { path: 'admin/users', component: UserListComponent },
];
