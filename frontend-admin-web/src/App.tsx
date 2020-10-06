import React, { useState } from "react";
import "./App.css";
import { Drawer, DrawerContent, DrawerHeader } from "@rmwc/drawer";
import {
  List,
  ListItem,
  ListItemGraphic,
  ListItemPrimaryText,
  ListItemSecondaryText,
  ListItemText,
} from "@rmwc/list";
import { Elevation } from "@rmwc/elevation";
import { ThemeProvider, Theme } from "@rmwc/theme";
import Logo from "./logo-white.svg";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useHistory,
} from "react-router-dom";
import DevicesPage from "./DevicesPage";
import AppDrawer from "./AppDrawer";
import UsersPage from "./UsersPage";
import ConnectPage from "./ConnectPage";
import LoginPage from "./LoginPage";
import { useLocation } from "react-router-dom";

const App: React.FC = () => {
  return (
    <Router>
      <ThemeProvider
        options={{
          primary: "rgba(57, 156, 248, 1)",
          secondary: "rgb(127, 193, 255)",
          onPrimary: "white",
          surface: "white",
        }}
        style={{
          height: "100%",
        }}
      >
        <div className="App">
          <AppContents />
        </div>
      </ThemeProvider>
    </Router>
  );
};

const AppContents: React.FC = () => {
  const location = useLocation();

  return (
    <>
      {location.pathname !== "/" && (
        <Elevation
          z={24}
          style={{
            height: "100%",
          }}
        >
          <AppDrawer />
        </Elevation>
      )}

      <Switch>
        <Route path="/devices">
          <DevicesPage />
        </Route>
        <Route path="/users">
          <UsersPage />
        </Route>
        <Route path="/connect">
          <ConnectPage />
        </Route>
        <Route path="/">
          <LoginPage />
        </Route>
      </Switch>
    </>
  );
};

export default App;
