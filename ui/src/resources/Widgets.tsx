import { 
    List,
    Datagrid,
    TextField,
    ReferenceField,
    Filter,
    DateField,
    EditButton,
    Edit,
    Create,
    SimpleForm,
    ReferenceInput,
    SelectInput,
    TextInput
} from 'react-admin';

const WidgetFilter = () => (
    <Filter>
        <TextInput label="Search" source="name" alwaysOn />
        <ReferenceInput label="Thing" source="thing_id" reference="things" allowEmpty>
            <SelectInput optionText="name" />
        </ReferenceInput>
    </Filter>
);

export const WidgetList = () => (
    <List filters={<WidgetFilter />}>
        <Datagrid rowClick="edit">
            <TextField source="name" />
            <TextField source="description" />
            <ReferenceField source="thing_id" reference="things">
                <TextField source="name" />
            </ReferenceField>
            <DateField source="created" showTime={true} />
            <DateField source="updated" showTime={true} />
            <EditButton />
        </Datagrid>
    </List>
);

export const WidgetEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="name" />
            <TextInput multiline source="description" />
            <ReferenceInput source="thing_id" reference="things">
                <SelectInput optionText="name" />
            </ReferenceInput>
        </SimpleForm>
    </Edit>
);


export const WidgetCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="name" />
            <TextInput multiline source="description" />
            <ReferenceInput source="thing_id" reference="things">
                <SelectInput optionText="name" />
            </ReferenceInput>
        </SimpleForm>
    </Create>
);