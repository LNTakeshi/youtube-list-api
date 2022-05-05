import React from 'react';
import { Navigate, useRoutes } from 'react-router-dom';
import Home from './Home';
import Redirect from './Redirect';
import Room from './Room';

export const Route = () => {
  const element = useRoutes([
    // These are the same as the props you provide to <Route>
    { path: process.env.PUBLIC_URL + '/', element: <Home /> },
    { path: process.env.PUBLIC_URL + '/room/:RoomId', element: <Room /> },
    {
      path: process.env.PUBLIC_URL + '/:RoomId',
      element: <Redirect />
    },
    // Not found routes work as you'd expect
    { path: '*', element: <Navigate to={process.env.PUBLIC_URL + '/'} /> }
  ]);

  // The returned element will render the entire element
  // hierarchy with all the appropriate context it needs
  return element;
};
