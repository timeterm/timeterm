import { Drawer, DrawerContent, DrawerHeader } from "@rmwc/drawer";
import Logo from "./logo-white.svg";
import { Theme, ThemeProvider } from "@rmwc/theme";
import {
  List,
  ListItem,
  ListItemGraphic,
  ListItemPrimaryText,
  ListItemSecondaryText,
  ListItemText,
} from "@rmwc/list";
import React from "react";
import { useHistory, useLocation } from "react-router-dom";

const AppDrawer: React.FC = () => {
  const history = useHistory();
  const location = useLocation();

  return (
    <Theme use={["primaryBg", "onPrimary"]} wrap>
      <Drawer>
        <DrawerHeader
          style={{
            marginTop: 16,
          }}
        >
          <img src={Logo} alt={"Timeterm Logo"} width={96} />
        </DrawerHeader>
        <DrawerContent>
          <ThemeProvider
            options={{
              primary: "rgb(0, 0, 0)",
            }}
            style={{
              display: "flex",
              flexDirection: "column",
              height: "100%",
            }}
          >
            <List style={{ flexGrow: 1 }}>
              <Theme use={["onPrimary"]} wrap>
                <ListItem
                  selected={location.pathname === "/devices"}
                  onClick={() => history.push("/devices")}
                >
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="tablet" />
                  </Theme>
                  Apparaten
                </ListItem>
              </Theme>
              <Theme use={["onPrimary"]} wrap>
                <ListItem
                  selected={location.pathname === "/users"}
                  onClick={() => history.push("/users")}
                >
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="group" />
                  </Theme>
                  Gebruikers
                </ListItem>
              </Theme>
              <Theme use={["onPrimary"]} wrap>
                <ListItem
                  selected={location.pathname === "/connect"}
                  onClick={() => history.push("/connect")}
                >
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="bluetooth_connected" />
                  </Theme>
                  Apparaat koppelen
                </ListItem>
              </Theme>
            </List>

            <List twoLine={true} style={{ paddingBottom: 0 }}>
              <Theme use={["onPrimary"]} wrap>
                <ListItem>
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="person" />
                  </Theme>
                  <ListItemText>
                    <ListItemPrimaryText>Admin</ListItemPrimaryText>
                    <Theme use={["onPrimary"]} wrap>
                      <ListItemSecondaryText>
                        admin@timeterm.nl
                      </ListItemSecondaryText>
                    </Theme>
                  </ListItemText>
                </ListItem>
              </Theme>
            </List>
            <List>
              <Theme use={["onPrimary"]} wrap>
                <ListItem>
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="logout" />
                  </Theme>
                  Uitloggen
                </ListItem>
              </Theme>
            </List>
          </ThemeProvider>
        </DrawerContent>
      </Drawer>
    </Theme>
  );
};

export default AppDrawer;
