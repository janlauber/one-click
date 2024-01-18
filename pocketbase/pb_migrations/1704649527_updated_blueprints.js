/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("vs5gr49hpah1g9q")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "ijfh1pyw",
    "name": "shared",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("vs5gr49hpah1g9q")

  // remove
  collection.schema.removeField("ijfh1pyw")

  return dao.saveCollection(collection)
})
