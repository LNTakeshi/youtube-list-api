import React, { createContext, Dispatch, useEffect, useState } from 'react';
import { Route } from './Route';
export type Theme = 'light' | 'dark';

const stylesheets = {
  light: process.env.PUBLIC_URL + '/antd.min.css',
  dark: process.env.PUBLIC_URL + '/antd.dark.min.css'
};

export const ColorCookie = 'color';
type ColorContextType = {
  theme: Theme;
  setTheme: Dispatch<React.SetStateAction<Theme>>;
};

const createStylesheetLink = (): HTMLLinkElement => {
  const link = document.createElement('link');
  link.rel = 'stylesheet';
  link.id = 'antd-stylesheet';
  document.head.appendChild(link);
  const preload = document.createElement('link');
  preload.rel = 'preload';
  preload.href = stylesheets['dark'];
  preload.as = 'style';
  document.head.appendChild(preload);
  const preload2 = document.createElement('link');
  preload2.rel = 'preload';
  preload2.href = stylesheets['light'];
  preload2.as = 'style';
  document.head.appendChild(preload2);
  return link;
};

export const getStylesheetLink = (): HTMLLinkElement =>
  document.head.querySelector('#antd-stylesheet') || createStylesheetLink();

const systemTheme = () =>
  window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light';

const getTheme = () =>
  (localStorage.getItem('theme') as Theme) || systemTheme();

export const updateTheme = (t: Theme) => {
  localStorage.setItem('theme', t);
  getStylesheetLink().href = stylesheets[t];
};

const defaultContext: ColorContextType = {
  theme: getTheme(),
  setTheme: () => {}
};

export const ColorContext = createContext<ColorContextType>(defaultContext);

type ForceReloadType = {
  load: boolean;
  setLoad: Dispatch<React.SetStateAction<boolean>>;
};

const defaultForceReload: ForceReloadType = {
  load: false,
  setLoad: () => {}
};
export const ForceReloadContext = createContext(defaultForceReload);
const App = () => {
  const [theme, setTheme] = useState<Theme>(getTheme());
  const [load, setLoad] = useState(false);
  updateTheme(theme);
  useEffect(() => {
    console.log('useEffect');
    const e = () => {
      console.log('setLoad');
      setLoad(true);
    };
    window.addEventListener('load', e);
    setTimeout(() => {
      if (!load) {
        setLoad(true);
      }
    }, 2000);
    return () => window.removeEventListener('load', e);
  }, [load]);

  return (
    <ForceReloadContext.Provider value={{ load, setLoad }}>
      <ColorContext.Provider value={{ theme, setTheme }}>
        {load && <Route />}
      </ColorContext.Provider>
    </ForceReloadContext.Provider>
  );
};

export default App;
