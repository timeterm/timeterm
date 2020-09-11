import React from "react";
import "./App.css";
import "@rmwc/drawer/styles";
import "@rmwc/list/styles";
import "@rmwc/button/styles";
import "@rmwc/checkbox/styles";
import "@rmwc/icon-button/styles";
import "@rmwc/icon/styles";
import "@rmwc/elevation/styles";
import "@rmwc/theme/styles";
import "@rmwc/data-table/styles";
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
import DevicesTable, { DeviceStatus } from "./DevicesTable";

function App() {
  return (
    <ThemeProvider
      options={{
        primary: "rgba(57, 156, 248, 1)",
        secondary: "rgb(127,193,255)",
        onPrimary: "white",
      }}
      style={{
        height: "100%",
      }}
    >
      <div className="App">
        <Elevation
          z={24}
          style={{
            height: "100%",
          }}
        >
          <Theme use={["primaryBg", "onPrimary"]} wrap>
            <Drawer>
              <DrawerHeader>
                <img src={Logo} alt={"Timeterm Logo"} width={96} />
              </DrawerHeader>
              <DrawerContent
                style={{
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <List style={{ flexGrow: 1 }}>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="tablet" />
                      </Theme>
                      Apparaten
                    </ListItem>
                  </Theme>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="group" />
                      </Theme>
                      Gebruikers
                    </ListItem>
                  </Theme>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="bluetooth_connected" />
                      </Theme>
                      Apparaat koppelen
                    </ListItem>
                  </Theme>
                </List>

                <List twoLine={true} style={{ paddingBottom: 0 }}>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="person" />
                      </Theme>
                      <ListItemText>
                        <ListItemPrimaryText>Admin</ListItemPrimaryText>
                        <Theme use={["primaryBg", "onPrimary"]} wrap>
                          <ListItemSecondaryText>
                            admin@timeterm.nl
                          </ListItemSecondaryText>
                        </Theme>
                      </ListItemText>
                    </ListItem>
                  </Theme>
                </List>
                <List>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="logout" />
                      </Theme>
                      Uitloggen
                    </ListItem>
                  </Theme>
                </List>
              </DrawerContent>
            </Drawer>
          </Theme>
        </Elevation>

        <Elevation
          z={16}
          style={{
            flexGrow: 1,
            margin: 16,
            borderRadius: 8,
          }}
        >
          <DevicesTable
            devices={[
              {
                name: "Mediatheek 1",
                status: DeviceStatus.Online,
              },
              {
                name: "Mediatheek 2",
                status: DeviceStatus.Offline,
              },
            ]}
          />
        </Elevation>
      </div>
    </ThemeProvider>
  );
}

export default App;
