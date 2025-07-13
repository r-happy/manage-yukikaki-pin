import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';

export const authGuard: CanActivateFn = (route, state) => {
  console.log('AuthGuard called for:', state.url);

  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    console.log('Token exists:', !!token);

    if (token) {
      return true;
    }

    console.log('No token, redirecting to signin');
    return inject(Router).createUrlTree(['/signin']);
  }

  // SSR時も未認証扱いでリダイレクト
  console.log('SSR environment, redirecting to signin');
  return inject(Router).createUrlTree(['/signin']);
};
