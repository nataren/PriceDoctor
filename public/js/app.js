API = "http://localhost:8080/@api/"
App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.ApplicationController = Ember.ArrayController.extend({
    queryParams: ['query', 'miles'],
    query: null,
    miles: null,
    
    queryField: Ember.computed.oneWay('query'),
    milesField: Ember.computed.oneWay('miles'),
    actions: {
        search: function() {
            this.set('query', this.get('queryField'));
            this.set('miles', 100 /* this.get('milesField') */);
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
        }
    },
    
    model: function(params) {
        if(!params.query) {
            return []; // no results;
        }
        var url = API + "healthproviders";
        return Ember.$.getJSON(url + "?address=" + params.query + "&" + "miles=" + params.miles).then(function(data) {
            return data.healthproviders
        });
    }
});
