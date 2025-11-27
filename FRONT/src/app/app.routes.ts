import { Routes } from '@angular/router';
import { HomepageComponent } from './homepage/homepage.component';
import { TourCreateComponent } from './tour/tour-create/tour-create.component';
import { BlogComponent } from './blog/blog.component';
import { TourListComponent } from './tour/tour-list/tour-list.component';
import { TourDetailComponent } from './tour/tour-detail/tour-detail.component';
import { LoginComponent } from './infrastructure/login/login.component'
import { RegisterComponent } from './infrastructure/register/register.component';
import { ProfileComponent } from './infrastructure/profile/profile/profile.component';
import { UserListComponent } from './auth/user-list/user-list.component';
import { ProfileFollowPageComponent } from './profile-follow/profile-follow-page/profile-follow-page.component';
import { ProfileRecommendedPageComponent } from './profile-follow/profile-recommended-page/profile-recommended-page.component';
import { LocationSimulatorComponent } from './locaiton-simulator/location-simulator/location-simulator.component';
import { ShoppingCartComponent } from './shopping-cart/shopping-cart/shopping-cart.component';

export const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomepageComponent },

    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },

    { path: 'tours', component: TourCreateComponent },
    { path: 'blog', component: BlogComponent },
    { path: 'tour-list', component: TourListComponent },
    { path: 'tours/:id', component: TourDetailComponent },
    { path: 'profile', component: ProfileComponent },
    { path: 'admin/users', component: UserListComponent },
    { path: 'profile/:id', component: ProfileFollowPageComponent },
    { path: 'recommendations', component: ProfileRecommendedPageComponent },
    { path: 'location-simulator', component: LocationSimulatorComponent },
    { path: 'shopping-cart', component: ShoppingCartComponent }
];