import { Card, CardContent, CardHeader } from '@mui/material';
import { Title } from 'react-admin';

const dashboard = () => (
    <Card>
        <Title title="Welcome to gorestapi" />
        <CardHeader title="Welcome to the gorestapi dashboard" />
        <CardContent>Select from the options on the left.</CardContent>
    </Card>
);

export default dashboard;