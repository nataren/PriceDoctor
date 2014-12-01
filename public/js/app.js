App = Ember.Application.create();

DS.Store.create({
    revision: 12,
    adapter: DS.RESTAdapter.create({
        namespace: 'api'
    })
});

App.Router.map(function() {
  // put your routes here
});

App.Kitten = DS.Model.extend({
    name: DS.attr('string'),
    picture: DS.attr('string')
});

App.IndexRoute = Ember.Route.extend({
    model: function() {
        return App.Kitten.find();
    }
});
