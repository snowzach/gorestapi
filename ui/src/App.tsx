import {
  Admin,
  Resource,
} from "react-admin";
import dataProvider from "./dataProvider";

import Dashboard from './resources/Dashboard';
import { ThingList, ThingEdit, ThingCreate } from './resources/Things';
import { WidgetList, WidgetEdit, WidgetCreate } from './resources/Widgets';

import LocalActivityIcon from '@mui/icons-material/LocalActivity'
import HotTubIcon from '@mui/icons-material/HotTub';

export const App = () => (
  <Admin disableTelemetry dataProvider={dataProvider} dashboard={Dashboard}>
    <Resource name="things" icon={LocalActivityIcon} list={ThingList} edit={ThingEdit} create={ThingCreate} />
    <Resource name="widgets" icon={HotTubIcon} list={WidgetList} edit={WidgetEdit} create={WidgetCreate} />
  </Admin>
);
