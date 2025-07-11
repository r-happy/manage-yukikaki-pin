import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';

export const authGuard: CanActivateFn = (route, state) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) return true;
    return inject(Router).createUrlTree(['/signin']);
  }
  // SSR時も未認証扱いでリダイレクト
  return inject(Router).createUrlTree(['/signin']);
};
