/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  collection.listRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.viewRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.createRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.updateRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.deleteRule = "@request.auth.id != \"\" && @request.auth.id = user.id"

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "ejjjrgsx",
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

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  collection.listRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.viewRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.createRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.updateRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"
  collection.deleteRule = "@request.auth.id != \"\" && @request.auth.id = serviceAccountName.id"

  // update
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
})
