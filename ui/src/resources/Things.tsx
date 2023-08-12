import { 
    List, 
    Edit, 
    Create,
    Filter,
    SimpleForm, 
    Datagrid,
    TextField, 
    DateField,
    EditButton, 
    TextInput,
    ReferenceInput,
    SelectInput,
    Pagination
} from 'react-admin';

const ListPagination = () => <Pagination rowsPerPageOptions={[10, 25, 50, 100]} />;

const listFilters = [
    <TextInput label="Search" source="q" alwaysOn />,
    <ReferenceInput label="Name" source="name" reference="things" allowEmpty>
        <SelectInput source="name" optionValue="name" optionText="name" />
    </ReferenceInput>
];

export const ThingList = () => (
    <List sort={{ field: 'id', order: 'DESC'}} pagination={<ListPagination />} filters={listFilters} >
       <Datagrid>
           <TextField source="name" />
           <TextField source="description" />
           <DateField source="created" showTime={true} />
           <DateField source="updated" showTime={true} />
           <EditButton />
        </Datagrid>
    </List>
);

export const ThingEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="id" />
            <TextInput source="name" />
            <TextInput multiline source="description" />
        </SimpleForm>
    </Edit>
);

export const ThingCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="name" />
            <TextInput multiline source="description" />
        </SimpleForm>
    </Create>
);
