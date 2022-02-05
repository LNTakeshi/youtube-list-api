import React, { useEffect } from 'react';
import { useMatch, useNavigate } from 'react-router-dom';

const Redirect = () => {
  const match = useMatch(process.env.PUBLIC_URL + '/:roomId');
  console.log(match);
  const navigate = useNavigate();
  useEffect(() => {
    if (match) {
      navigate(process.env.PUBLIC_URL + '/room/' + match.params.roomId);
    } else {
      navigate(process.env.PUBLIC_URL + '/');
    }
  });

  return <></>;
};
export default Redirect;
