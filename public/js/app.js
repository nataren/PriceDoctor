API = "http://localhost:8080/@api/"
App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.ApplicationController = Ember.ArrayController.extend({
    queryParams: ['query'],
    query: null,
    
    queryField: Ember.computed.oneWay('query'),
    actions: {
        search: function() {
            this.set('query', this.get('queryField'));
        }
    }
});

App.ApplicationRoute = Ember.Route.extend({
    queryParams: {
        query: {
            // Opt into full transition
            refreshModel: true
        }
    },
    
    model: function(params) {
        if(!params.query) {
            return []; // no results;
        }
        var url = API + "healthproviders";
        return Ember.$.getJSON(url + "?address=" + params.query).then(function(data) {
            return data.healthproviders
        });
    }
});
