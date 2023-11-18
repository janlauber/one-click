/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  collection.listRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.viewRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.createRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.updateRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.deleteRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"

  // remove
  collection.schema.removeField("dlyhvyrc")

  // remove
  collection.schema.removeField("5ldsmkui")

  // remove
  collection.schema.removeField("bwfhbyer")

  // remove
  collection.schema.removeField("ydvamb9k")

  // remove
  collection.schema.removeField("kktrixuw")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "ejjjrgsx",
    "name": "serviceAccountName",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "_pb_users_auth_",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  collection.listRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.viewRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.createRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.updateRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.deleteRule = "@request.auth.id != \"\" && @request.auth.id = user.id"

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "dlyhvyrc",
    "name": "user",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "_pb_users_auth_",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "5ldsmkui",
    "name": "values",
    "type": "text",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "bwfhbyer",
    "name": "version",
    "type": "text",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "ydvamb9k",
    "name": "startDate",
    "type": "date",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "min": "",
      "max": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "kktrixuw",
    "name": "endDate",
    "type": "date",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": "",
      "max": ""
    }
  }))

  // remove
  collection.schema.removeField("ejjjrgsx")

  return dao.saveCollection(collection)
})
