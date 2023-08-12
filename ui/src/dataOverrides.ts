

type SearchOverrides = { [index: string]: Search };
type SortOverrides = { [index: string]: string };

type Search = {
    operator?: string
};

type Override = {
    path?: string
    idField?: string
    search?: SearchOverrides
    sort?: SortOverrides
};

type Overrides = { [index: string]: Override };

const overrides: Overrides = {
    whateverResource: {
        path: 'whatever/path',
        idField: 'whatever_id',
        search: {
            'name': { operator: "=~~" },
            'whatever_field': { operator: "=~~" },
        }
    },
    portalConfig: {
        path: 'portal/config',
    },
    strategy: {
        idField: 'source',
    },
    quotes: {
        path: 'quotes',
        search: {
            'list': { operator: "=~~" },
            'address': { operator: "=~~" },
            'comment': { operator: "=~~" },
        },
        sort: {
            'createdAt.value': '-createdAt.value'
        }
    },
    rfqConfig: {
        path: 'rfq/config',
        search: {
            'id': { operator: "=~~" },
        }
    },
    rfqConfigSafeties: {
        path: 'rfq/config/safeties',
        search: {
            'id': { operator: "=~~" },
        }
    },
    rfqConfigFees: {
        path: 'rfq/config/fees',
        search: {
            'id': { operator: "=~~" },
        }
    },
    assetClasses: {
        path: 'assets/class',
        search: {
            'name': { operator: "=~~" },
            'description': { operator: "=~~" },
        }
    },
    portalUsers: {
        path: 'portal/users',
        search: {
            'name': { operator: "=~~" },
            'description': { operator: "=~~" },
        }
    },
    hrAssets: {
        path: 'hiddenroad/assets',
        idField: 'symbol',
    },
    hrCounterparties: {
        path: 'hiddenroad/counterparties',
        idField: 'counterparty',
    },
    hrAccountBalances: {
        path: 'hiddenroad/accounts/balances',
        idField: 'account.id',
    },
    hrInstruments: {
        path: 'hiddenroad/instruments',
        idField: 'symbol',
    },
    hrMarkets: {
        path: 'hiddenroad/markets',
        idField: 'symbol',
    },
    hrTrades: {
        path: 'hiddenroad/trades',
        idField: 'trade_id',
    },
}

export default overrides;