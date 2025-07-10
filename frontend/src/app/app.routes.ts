import { Routes } from '@angular/router';
import { Home } from './home/home';
import { Signin } from './signin/signin';
import { Signup } from './signup/signup';

export const routes: Routes = [
  {
    path: '',
    title: 'Home',
    component: Home,
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
