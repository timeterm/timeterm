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
import { IconButton } from "@rmwc/icon-button";
import { Checkbox } from "@rmwc/checkbox";
import {
  DataTable,
  DataTableBody,
  DataTableCell,
  DataTableContent,
  DataTableHead,
  DataTableHeadCell,
  DataTableRow,
} from "@rmwc/data-table";
import Logo from "./logo-white.svg";
import { Icon } from "@rmwc/icon";

function App() {
  return (
    <ThemeProvider
      options={{
        primary: "rgba(57, 156, 248, 1)",
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
          <DataTable
            style={{
              width: "100%",
              height: "100%",
            }}
          >
            <DataTableContent>
              <DataTableHead>
                <DataTableRow>
                  <DataTableHeadCell hasFormControl style={{whiteSpace: "nowrap"}}>
                    <Checkbox />
                  </DataTableHeadCell>
                  <DataTableHeadCell style={{width: "54.5%"}}>Naam</DataTableHeadCell>
                  <DataTableHeadCell style={{width: "54.5%"}}>Status</DataTableHeadCell>
                </DataTableRow>
              </DataTableHead>
              <DataTableBody>
                <DataTableRow>
                    <DataTableCell hasFormControl style={{whiteSpace: "nowrap"}}>
                        <Checkbox />
                    </DataTableCell> 

                  <DataTableCell style={{width: "54.5%"}}>Mediatheek 1</DataTableCell>
                  <DataTableCell style={{display: "inline-flex", alignItems: "center", width: "54.5%"}}>
                    <Icon
                      icon="check_circle"
                      style={{
                        color: "#4ECD6A",
                      }}
                    />&nbsp;
                    Online
                  </DataTableCell>
                </DataTableRow>
              </DataTableBody>
            </DataTableContent>
          </DataTable>
        </Elevation>
      </div>
    </ThemeProvider>
  );
}

export default App;
