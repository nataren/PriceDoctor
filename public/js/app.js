API = "http://localhost:8080/@api/"
App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.ApplicationController = Ember.ArrayController.extend({
    queryParams: ['query', 'miles', 'procedure'],
    query: null,
    miles: null,
    procedure: null,
    distance: [5, 10, 25, 50, 100, 250, 500],

    queryField: Ember.computed.oneWay('query'),
    milesField: Ember.computed.oneWay('miles'),
    procedureField: Ember.computed.oneWay('procedure'),
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
        miles: {
            refreshModel: true
        }
    },

    model: function(params) {
        if(!params.query) {
            return []; // no results;
        }
        var url = API + "healthproviders";
        return Ember.$.getJSON(url + "?address=" + params.query + "&" + "miles=" + params.miles + "&procedure=" + params.procedure).then(function(data) {
            return data
        });
    }
});
