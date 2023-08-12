import { DataProvider } from 'ra-core';
import { each, isObject, isArray } from 'lodash';
import dataOverrides from './dataOverrides';

import { fetchUtils, HttpError } from 'react-admin';

export const apiUrl = import.meta.env.VITE_API_URL;
const httpClient = fetchUtils.fetchJson;

// If we want to also support authentication
// const httpClient = (url: string, options = {}) => {
//     const token = localStorage.getItem('token');
//     let user = {};
//     if (token) {
//         user = { token: `Bearer ${token}`, authenticated: !!token };
//     }
//    return fetchUtils.fetchJson(url, {...options, user});
// }

const api: DataProvider = {

    // getList returns a list of resources
    getList: (resource, params) => {

        // Build the filter
        let query = buildFilter(resource, params);

        // Build Sort
        if(params.sort) {
            const { field, order } = params.sort;
            const sortName = getObjectField(dataOverrides, resource, "sort", field)
            if(sortName) {
                query += `&sort=${order === 'DESC' ? '-' : ''}${sortName}`
            } else {
                query += `&sort=${order === 'DESC' ? '-' : ''}${field}`
            }
            
        }

        // Build Pagination
        if(params.pagination) {
            const { page, perPage } = params.pagination;
            query += `&offset=${(page-1)*perPage}&limit=${perPage}`;
        }

        // Build URL
        const url = `${apiUrl}/${pathByResource(resource)}?`+query.substring(1);

        return httpClient(url).then(({ headers, json }) => {
            // If the id field in the record is overridden it must be transformed.
            const idField = getObjectField(dataOverrides, resource, "idField")
            if (idField)  {
                json.results.forEach((row: any, index: number) => {
                    row.id = row[idField];
                    json.results[index] = row;
                });
            }

            // If pagination is enabled, there are results and no/zero count
            // we need to provide the pageInfo object
            if(params.pagination && json.results.length && !json.count) {
                return {
                    data: json.results,
                    pageInfo: {
                        hasPreviousPage: (params.pagination.page > 1),
                        hasNextPage: json.results.length == params.pagination.perPage,
                    }
                };
            }

            return {
                data: json.results,
                total: json.count,
            }
        });
    },

    // getOne gets one record by id
    getOne: (resource, params) => {
        return httpClient(`${apiUrl}/${pathByResource(resource)}/${params.id}`).then(({ json }) => {
            // If the id field in the record is overridden it must be transformed.
            const idField = getObjectField(dataOverrides, resource, "idField")
            if (idField)  {
                json.id = json[idField];
            }
            return {
                data: json,
            }
        })

    },

    // getMany fetches many records by IDs
    getMany: (resource, params) => {
        const query = `id=(${params.ids.join(',')})`;
        const url = `${apiUrl}/${pathByResource(resource)}?${query}`;

        return httpClient(url).then(({ json }) => {
            // If the id field in the record is overridden it must be transformed.
            const idField = getObjectField(dataOverrides, resource, "idField")
            if (idField)  {
                json.results.forEach((row: any, index: number) => {
                    row.id = row[idField];
                    json.results[index] = row;
                });
            }

            return ({
                data: json.results 
            });
        });
    },

    // getManyReference returns a list of resource related to another resource (all comments for a post)
    getManyReference: (resource: string, params: any) => {
        
        let query = '';
        let filter = flattenObject(params.filter);

        //  Build Filter
        for(let field in filter) {
            query += `&${field}=`
            let value = filter[field];
            if(typeof value == 'object') {
                query += `(${value.join(',')})`;
            } else {
                query += `"${value}"`;
            }
        }

        query += `&${params.target}=`
        if(typeof params.id == 'object') {
            query += `(${params.id.join(',')})`;
        } else {
            query += `"${params.id}"`;
        }

        // Build Sort
        const { field, order } = params.sort;
        if(order) {
            query += `&sort=${order === 'DESC' ? '-' : ''}${field}`
        }

        // Build Pagination
        const { page, perPage } = params.pagination;
        query += `&offset=${(page-1)*perPage}&limit=${perPage}`;

        // Build URL
        const url = `${apiUrl}/${pathByResource(resource)}?`+query.substring(1);

        return httpClient(url).then(({ headers, json }) => ({
            data: json.results,
            total: json.count,
        }));
    },

    // update updates one resource
    update: (resource, params) =>
        httpClient(`${apiUrl}/${pathByResource(resource)}`, {
            method: 'POST',
            body: JSON.stringify(params.data),
        }).then(({ json }) => {
            // If the id field in the record is overridden it must be transformed.
            const idField = getObjectField(dataOverrides, resource, "idField")
            if (idField)  {
                json.id = json[idField];
            }
            return {
                data: json,
            }
        }).catch((err) => {
            if(err.body.error) {
                return Promise.reject(new HttpError(err.body.error, err.status));
            }
            return Promise.reject(err);
        }),

    // updateMany updates many resources
    updateMany: (resource, params) => {
        console.log("Not supported...");
        return Promise.resolve({data:[]});
    },

    // create creates on resource
    create: (resource, params) =>
        httpClient(`${apiUrl}/${pathByResource(resource)}`, {
            method: 'POST',
            body: JSON.stringify(params.data),
        }).then(({ json }) => {
            // If the id field in the record is overridden it must be transformed.
            const idField = getObjectField(dataOverrides, resource, "idField")
            if (idField)  {
                json.id = json[idField];
            }
            return {
                data: json,
            }
        }).catch((err) => {
            if(err.body.error) {
                return Promise.reject(new HttpError(err.body.error, err.status));
            }
            return Promise.reject(err);
        }),

    // delete a resource
    delete: (resource, params) =>
        httpClient(`${apiUrl}/${pathByResource(resource)}/${params.id}`, {
            method: 'DELETE',
        }).then(({json}) => ({
            data: { id: params.id, ...json },
        })).catch((err) => {
            if(err.body.error) {
                return Promise.reject(new HttpError(err.body.error, err.status));
            }
            return Promise.reject(err);
        }),

    // delete many resources
    deleteMany: (resource, params) => {
        // const query = {
        //     filter: JSON.stringify({ id: params.ids}),
        // };
        // return httpClient(`${apiUrl}/${pathByResource(resource)}?${stringify(query)}`, {
        //     method: 'DELETE',
        // }).then(({ json }) => ({ data: json }));
        console.log("Not supported...");
        return Promise.resolve({data:[]});
    },

    // get fetches anything from the API
    get: (path: string) => {
        const url = `${apiUrl}/${path}`;
        return httpClient(url).then(({ json }) => (json));;
    },

};

export default api;

// getObjectField traverses an object by field names in args and returns the field
const getObjectField = (obj: any, ...args: string[]): any => {
    return args.reduce((obj, level) => obj && obj[level], obj)
}

// pathByResource gets the path for a resource taking into consideration and overrides
export const pathByResource = (resource: string): string => {
    if (resource in dataOverrides) {
        let override = dataOverrides[resource];
        if (override && override.path) {
            return override.path;
        }
    }
    return resource;
}

// flattenObject turns an object into a single level object with keys of period delimited values.
export const flattenObject = (obj: any): {[key: string]: any} => {
    var nobj: {[key: string]: any} = {};
    each(obj, function(val: any, key: any){
        // ensure is JSON key-value map, not array
        if (isObject(val) && !isArray(val)) {
            // union the returned result by concat all keys
            var strip = flattenObject(val)
            each(strip, function(v,k){
                nobj[key+'.'+k] = v
            })
        } else {
            nobj[key] = val
        }
    })
    return nobj
}

// buildFilter will build the queryp compatible filter from the search options
// taking into account any overrides.
const buildFilter = (resource: string, params: any):string => {

    const searchOverrides = getObjectField(dataOverrides, resource, "search")

    // If search is a function, call it to get the filter
    if (searchOverrides instanceof Function) {
        return searchOverrides(params);
    }

    let ret = '';

    // Handle the various forms of filter.option
    if(params.filter && params.filter.option) {
        const options = params.filter.option
        if(options instanceof Object ) {
            if(options instanceof Array) {
                options.forEach((option) => {
                    ret += `&option=${option}`;
                });
            } else {
                for(let option in options){
                    ret += `&option[${option}]=${options[option]}`;
                }
            }
        } else {
            ret += `&option=${options}`
        }
        delete params.filter.option;
    }

    const filter = flattenObject(params.filter);
      
    // If search values have been overridden
    if (searchOverrides) {
        for(let field in filter) {
            const search = searchOverrides[field];
            const value = filter[field];
            if(value === undefined) continue
            if (search) {
                if (search.field) {
                    ret += `&${encodeURIComponent(search.field)}`
                } else {
                    ret += `&${encodeURIComponent(field)}`
                }
                if (search.operator) {
                    ret += search.operator
                } else {
                    ret += '='
                }
            } else {
                ret += `&${field}=`
            }
            if(typeof value == 'object') {
                ret += `("${encodeURIComponent(value.join('","'))}")`;
            } else if(typeof value == 'number' || typeof value == 'boolean') {
                ret += value.toString();
            } else {
                if(value.charAt(0) === '(' && value.charAt(value.length-1) === ')') {
                    ret += `${encodeURIComponent(value)}`;
                } else {
                    ret += `"${encodeURIComponent(value)}"`;
                }
            }
        }
    } else {
        for(let field in filter) {
            const value = filter[field];
            if(value === undefined) continue
            ret += `&${field}=`
            if(typeof value == 'object') {
                ret += `(${encodeURIComponent(value.join(','))})`;
            } else if(typeof value == 'number') {
                ret += value.toString();
            } else {
                if(value.charAt(0) === '(' && value.charAt(value.length-1) === ')') {
                    ret += `${encodeURIComponent(value)}`;
                } else {
                    ret += `"${encodeURIComponent(value)}"`;
                }
            }
        }
    }

    return ret;
}
