import React from "react";
import "./App.css";
import { Elevation } from "@rmwc/elevation";
import { ThemeProvider } from "@rmwc/theme";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";
import DevicesPage from "./DevicesPage";
import AppDrawer from "./AppDrawer";
import UsersPage from "./UsersPage";
import LoginPage from "./LoginPage";
import { useLocation } from "react-router-dom";
import LoginDonePage from "./LoginDonePage";
import Cookies from "universal-cookie";
import { QueryCache, ReactQueryCacheProvider } from "react-query";

const App: React.FC = () => {
  return (
    <Router>
      <ThemeProvider
        options={{
          primary: "rgb(57,156,248)",
          secondary: "rgb(127,193,255)",
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

const queryCache = new QueryCache();

const AppContents: React.FC = () => {
  const location = useLocation();
  const session = new Cookies().get("ttsess");
  const loggedIn = !!session;

  return (
    <ReactQueryCacheProvider queryCache={queryCache}>
      {!["/", "/login/done"].includes(location.pathname) && (
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
        <Route path={"/login/done"}>
          <LoginDonePage />
        </Route>
        <Route exact path="/">
          <LoginPage />
        </Route>

        {!loggedIn && <Redirect to={"/"} />}

        <Route path="/devices">
          <DevicesPage />
        </Route>
        <Route path="/users">
          <UsersPage />
        </Route>
      </Switch>
    </ReactQueryCacheProvider>
  );
};

export default App;
