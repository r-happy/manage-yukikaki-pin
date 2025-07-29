import { Routes } from '@angular/router';
import { Signin } from './signin/signin';
import { Signup } from './signup/signup';
import { Dashboard } from './dashboard/dashboard';
import { authGuard } from './guard/auth-guard';
import { Profile } from './profile/profile';
import { Map } from './map/map';
import { Pin } from './pin/pin';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/dashboard',
    pathMatch: 'full',
  },
  {
    path: 'dashboard',
    title: 'Dashboard',
    component: Dashboard,
    canActivate: [authGuard],
    children: [
      {
        path: '',
        redirectTo: 'map',
        pathMatch: 'full',
      },
      {
        path: 'map',
        title: 'Map',
        component: Map,
      },
      {
        path: 'pin',
        title: 'Pin',
        component: Pin,
      },
      {
        path: 'profile',
        title: 'Profile',
        component: Profile,
      },
    ],
  },
  {
    path: 'signin',
    title: 'Sign In',
    component: Signin,
  },
  {
    path: 'signup',
    title: 'Sign Up',
    component: Signup,
  },
];
