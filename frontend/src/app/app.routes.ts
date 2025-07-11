import { Routes } from '@angular/router';
import { Home } from './home/home';
import { Signin } from './signin/signin';
import { Signup } from './signup/signup';
import { Dashboard } from './dashboard/dashboard';
import { authGuard } from './guard/auth-guard';
import { Layout } from './layout/layout';
import { Profile } from './profile/profile';

export const routes: Routes = [
  {
    path: '',
    component: Layout,
    canActivate: [authGuard],
    children: [
      {
        path: '',
        title: 'Home',
        component: Home,
      },
      {
        path: 'dashboard',
        title: 'Dashboard',
        component: Dashboard,
        children: [{ path: 'profile', title: 'Profile', component: Profile }],
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
