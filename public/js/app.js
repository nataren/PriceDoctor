API = "@api/"
App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.ApplicationController = Ember.ArrayController.extend({
    queryParams: ['query', 'miles', 'procedure'],
    query: null,
    miles: null,
    procedure: null,
    distance: [5, 10, 25, 50, 100, 250, 500],
    sortBy: ["price"],

    queryField: Ember.computed.oneWay('query'),
    milesField: Ember.computed.oneWay('miles'),
    procedureField: Ember.computed.oneWay('procedure'),
    sortByField: Ember.computed.oneWay('sortBy'),
    actions: {
        search: function() {
            this.set('query', this.get('queryField'));
            this.set('miles', this.get('milesField'));
            this.set('procedure', this.get('procedureField'));
        }
    }
});

App.ApplicationRoute = Ember.Route.extend({
    queryParams: {
        query: {
            // Opt into full transition
            refreshModel: true
        },
        miles: {
            refreshModel: true
        },
        procedure: {
            refreshModel: true
        },
        sortBy: {
            refreshModel: true
        }
    },

    model: function(params) {
        if(!params.query) {
            return []; // no results;
        }
        return Ember.$.getJSON(API + "healthproviders" + "?address=" + params.query + "&" + "miles=" + params.miles + "&procedure=" + params.procedure + "&sortby=" + params.sortBy).then(function(data) {
            return data
        });
    }
});
